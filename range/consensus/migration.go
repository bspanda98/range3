// Copyright 2019 The Range Core Authors
// This file is part of the Range Core library.
//
// The Range Core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Range Core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Range Core library. If not, see <http://www.gnu.org/licenses/>.

package consensus

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math/big"
	"os"
	"strings"

	"range/core/gen3/accounts/abi"
	"range/core/gen3/common"
	"range/core/gen3/consensus"
	"range/core/gen3/core"
	"range/core/gen3/core/state"
	"range/core/gen3/core/types"
	"range/core/gen3/log"
	"range/core/gen3/params"
	"range/core/gen3/rlp"

	"github.com/shengdoushi/base58"

	energi_abi "range/core/gen3/energi/abi"
	energi_params "range/core/gen3/energi/params"
)

const (
	gasPerMigrationEntry uint64 = 100000
)

func (e *Range) finalizeMigration(
	chain ChainReader,
	header *types.Header,
	statedb *state.StateDB,
	txs types.Transactions,
) error {
	if !header.IsGen2Migration() {
		return nil
	}

	// One migration and one block reward
	if len(txs) != 2 {
		err := errors.New("Wrong number of migration block txs")
		log.Error("Failed to finalize migration", "err", err)
		return err
	}

	migration_abi, err := abi.JSON(strings.NewReader(energi_abi.Gen2MigrationABI))
	if err != nil {
		return err
	}

	callData, err := migration_abi.Pack("totalAmount")
	if err != nil {
		log.Error("Failed to prepare totalAmount() call", "err", err)
		return err
	}

	// totalAmount()
	msg := types.NewMessage(
		e.systemFaucet,
		&energi_params.Range_MigrationContract,
		0,
		common.Big0,
		e.callGas,
		common.Big0,
		callData,
		false,
	)
	rev_id := statedb.Snapshot()
	evm := e.createEVM(msg, chain, header, statedb)
	gp := core.GasPool(e.callGas)
	output, _, _, err := core.ApplyMessage(evm, msg, &gp)
	statedb.RevertToSnapshot(rev_id)
	if err != nil {
		log.Error("Failed in totalAmount() call", "err", err)
		return err
	}

	//
	totalAmount := big.NewInt(0)
	err = migration_abi.Unpack(&totalAmount, "totalAmount", output)
	if err != nil {
		log.Error("Failed to unpack totalAmount() call", "err", err)
		return err
	}

	statedb.SetBalance(energi_params.Range_MigrationContract, totalAmount)
	log.Warn("Setting Migration contract balance", "amount", totalAmount)

	return nil
}

func MigrationTx(
	signer types.Signer,
	header *types.Header,
	migration_file string,
	engine consensus.Engine,
) (res *types.Transaction) {
	file, err := os.Open(migration_file)
	if err != nil {
		log.Error("Failed to open snapshot", "err", err)
		return nil
	}
	defer file.Close()

	snapshot, err := parseSnapshot(file)
	if err != nil {
		log.Error("Failed to parse snapshot", "err", err)
		return nil
	}

	return migrationTx(signer, header, snapshot, engine)
}

func migrationTx(
	signer types.Signer,
	header *types.Header,
	snapshot *snapshot,
	engine consensus.Engine,
) (res *types.Transaction) {
	e, ok := engine.(*Range)
	if !ok {
		log.Error("Not Range consensus engine")
		return nil
	}

	owners, amounts, blacklist := createSnapshotParams(snapshot)
	if owners == nil || amounts == nil || blacklist == nil {
		log.Error("Failed to create arguments")
		return nil
	}

	migration_abi, err := abi.JSON(strings.NewReader(energi_abi.Gen2MigrationABI))
	if err != nil {
		panic(err)
	}

	callData, err := migration_abi.Pack("setSnapshot", owners, amounts, blacklist)
	if err != nil {
		panic(err)
	}

	gasLimit := gasPerMigrationEntry * uint64(len(owners))
	header.GasLimit = gasLimit
	header.Extra, err = rlp.EncodeToBytes([]interface{}{
		uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
		"range3",
		snapshot.Hash,
	})
	if err != nil {
		panic(err)
	}

	res = types.NewTransaction(
		uint64(0), // it should be the first transaction
		energi_params.Range_MigrationContract,
		common.Big0,
		gasLimit,
		common.Big0,
		callData,
	)

	if e.signerFn == nil {
		log.Error("Signer is not set")
		return nil
	}

	if e.config == nil {
		log.Error("Engine config is not set")
		return nil
	}

	tx_hash := signer.Hash(res)
	tx_sig, err := e.signerFn(e.config.MigrationSigner, tx_hash.Bytes())
	if err != nil {
		log.Error("Failed to sign migration tx")
		return nil
	}

	res, err = res.WithSignature(signer, tx_sig)
	if err != nil {
		log.Error("Failed to pack migration tx")
		return nil
	}
	return
}

func createSnapshotParams(ss *snapshot) (
	owners []common.Address,
	amounts []*big.Int,
	blacklist []common.Address,
) {
	owners = make([]common.Address, len(ss.Txouts))
	amounts = make([]*big.Int, len(ss.Txouts))
	blacklist = make([]common.Address, len(ss.Blacklist))

	// NOTE: Gen 2 precision is 8, but Gen 3 is 18
	multiplier := big.NewInt(1e10)

	for i, info := range ss.Txouts {
		owner, err := base58.Decode(info.Owner, base58.BitcoinAlphabet)

		if err != nil {
			log.Error("Failed to decode address", "err", err, "address", info.Owner)
			return nil, nil, nil
		}

		owner = owner[1 : len(owner)-4]
		owners[i] = common.BytesToAddress(owner)
		amounts[i] = new(big.Int).Mul(info.Amount, multiplier)
	}

	for i, blo := range ss.Blacklist {
		owner, err := base58.Decode(blo, base58.BitcoinAlphabet)

		if err != nil {
			log.Error("Failed to decode address", "err", err, "address", blo)
			return nil, nil, nil
		}

		owner = owner[1 : len(owner)-4]
		blacklist[i] = common.BytesToAddress(owner)
	}

	return
}

func parseSnapshot(reader io.Reader) (*snapshot, error) {
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	ret := &snapshot{}
	err := dec.Decode(ret)
	return ret, err
}

type snapshotItem struct {
	Owner  string   `json:"owner"`
	Amount *big.Int `json:"amount"`
	Atype  string   `json:"type"`
}

type snapshot struct {
	Txouts    []snapshotItem `json:"snapshot_utxos"`
	Blacklist []string       `json:"snapshot_blacklist"`
	Hash      string         `json:"snapshot_hash"`
}

func ValidateMigration(
	block *types.Block,
	migration_file string,
) bool {
	file, err := os.Open(migration_file)
	if err != nil {
		log.Error("Failed to open snapshot", "err", err)
		return false
	}
	defer file.Close()

	snapshot, err := parseSnapshot(file)
	if err != nil {
		log.Error("Failed to parse snapshot", "err", err)
		return false
	}

	owners, amounts, blacklist := createSnapshotParams(snapshot)
	if owners == nil || amounts == nil || blacklist == nil {
		log.Error("Failed to create arguments")
		return false
	}

	migration_abi, err := abi.JSON(strings.NewReader(energi_abi.Gen2MigrationABI))
	if err != nil {
		panic(err)
	}

	callData, err := migration_abi.Pack("setSnapshot", owners, amounts, blacklist)
	if err != nil {
		panic(err)
	}

	txs := block.Transactions()
	if len(txs) != 2 {
		log.Error("Invalid transaction count")
		return false
	}

	if !bytes.Equal(txs[0].Data(), callData) {
		log.Error("Migration transaction data mismatch")
		return false
	}

	return true
}
