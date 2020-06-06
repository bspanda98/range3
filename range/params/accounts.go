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

package params

import (
	"math/big"

	"range/core/gen3/common"
)

var (
	Range_BlockReward        = common.BigToAddress(big.NewInt(0x300))
	Range_Treasury           = common.BigToAddress(big.NewInt(0x301))
	Range_MasternodeRegistry = common.BigToAddress(big.NewInt(0x302))
	Range_StakerReward       = common.BigToAddress(big.NewInt(0x303))
	Range_BackboneReward     = common.BigToAddress(big.NewInt(0x304))
	Range_SporkRegistry      = common.BigToAddress(big.NewInt(0x305))
	Range_CheckpointRegistry = common.BigToAddress(big.NewInt(0x306))
	Range_BlacklistRegistry  = common.BigToAddress(big.NewInt(0x307))
	Range_MigrationContract  = common.BigToAddress(big.NewInt(0x308))
	Range_MasternodeToken    = common.BigToAddress(big.NewInt(0x309))
	Range_Blacklist          = common.BigToAddress(big.NewInt(0x30A))
	Range_Whitelist          = common.BigToAddress(big.NewInt(0x30B))
	Range_MasternodeList     = common.BigToAddress(big.NewInt(0x30C))

	Range_BlockRewardV1        = common.BigToAddress(big.NewInt(0x310))
	Range_TreasuryV1           = common.BigToAddress(big.NewInt(0x311))
	Range_MasternodeRegistryV1 = common.BigToAddress(big.NewInt(0x312))
	Range_StakerRewardV1       = common.BigToAddress(big.NewInt(0x313))
	Range_BackboneRewardV1     = common.BigToAddress(big.NewInt(0x314))
	Range_SporkRegistryV1      = common.BigToAddress(big.NewInt(0x315))
	Range_CheckpointRegistryV1 = common.BigToAddress(big.NewInt(0x316))
	Range_BlacklistRegistryV1  = common.BigToAddress(big.NewInt(0x317))
	Range_CompensationFundV1   = common.BigToAddress(big.NewInt(0x318))
	Range_MasternodeTokenV1    = common.BigToAddress(big.NewInt(0x319))

	Range_SystemFaucet = common.BigToAddress(big.NewInt(0x320))
	Range_Ephemeral    = common.HexToAddress("0x457068656d6572616c")

	// NOTE: this is NOT very safe, but it optimizes significantly
	Storage_ProxyImpl = common.BigToHash(big.NewInt(0x01))
)
