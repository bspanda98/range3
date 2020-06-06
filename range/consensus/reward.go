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
	"math/big"

	"range/core/gen3/common"
	"range/core/gen3/common/math"
	"range/core/gen3/core"
	"range/core/gen3/core/state"
	"range/core/gen3/core/types"
	"range/core/gen3/log"

	energi_params "range/core/gen3/energi/params"
)

var (
	BigBalance = new(big.Int).Div(math.MaxBig256, big.NewInt(2))
)

func (e *Range) processBlockRewards(
	chain ChainReader,
	header *types.Header,
	statedb *state.StateDB,
	txs types.Transactions,
	receipts types.Receipts,
) (types.Transactions, types.Receipts, error) {
	systemFaucet := e.systemFaucet

	// Temporary balance setup & clean up
	statedb.SetBalance(systemFaucet, BigBalance)
	defer statedb.SetBalance(systemFaucet, common.Big0)

	// Common get reward call
	getRewardData, err := e.rewardAbi.Pack("getReward", header.Number)
	if err != nil {
		log.Error("Fail to prepare getReward() call", "err", err)
		return nil, nil, err
	}

	rewardData, err := e.rewardAbi.Pack("reward")
	if err != nil {
		log.Error("Fail to prepare reward() call", "err", err)
		return nil, nil, err
	}

	// GetReward()
	//====================================
	msg := types.NewMessage(
		systemFaucet,
		&energi_params.Range_BlockReward,
		0,
		common.Big0,
		e.callGas,
		common.Big0,
		getRewardData,
		false,
	)
	rev_id := statedb.Snapshot()
	evm := e.createEVM(msg, chain, header, statedb)
	gp := core.GasPool(msg.Gas())
	output, gas1, _, err := core.ApplyMessage(evm, msg, &gp)
	statedb.RevertToSnapshot(rev_id)
	if err != nil {
		log.Error("Failed in getReward() call", "err", err)
		return nil, nil, err
	}

	//
	total_reward := big.NewInt(0)
	err = e.rewardAbi.Unpack(&total_reward, "getReward", output)
	if err != nil {
		log.Error("Failed to unpack getReward() call", "err", err)
		return nil, nil, err
	}

	// Reward
	//====================================
	tx := types.NewTransaction(
		statedb.GetNonce(systemFaucet),
		energi_params.Range_BlockReward,
		total_reward,
		e.xferGas,
		common.Big0,
		rewardData)
	tx = tx.WithConsensusSender(systemFaucet)

	statedb.Prepare(tx.Hash(), header.Hash(), len(txs))

	msg, err = tx.AsMessage(&ConsensusSigner{})
	if err != nil {
		log.Error("Failed in BlockReward AsMessage()", "err", err)
		return nil, nil, err
	}

	evm = e.createEVM(msg, chain, header, statedb)
	gp = core.GasPool(msg.Gas())
	_, gas2, failed, err := core.ApplyMessage(evm, msg, &gp)
	if err != nil {
		log.Error("Failed in reward() call", "err", err)
		return nil, nil, err
	}

	// NOTE: it should be Byzantium finalization here...
	root := statedb.IntermediateRoot(chain.Config().IsEIP158(header.Number))
	receipt := types.NewReceipt(root.Bytes(), failed, header.GasUsed)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = gas2
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	log.Trace("Block reward", "reward", total_reward, "gas", gas1+gas2)

	return append(txs, tx), append(receipts, receipt), nil
}
