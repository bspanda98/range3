// Copyright 2019 The Range Core Authors
// This file is part of Range Core.
//
// Range Core is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Range Core is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Range Core. If not, see <http://www.gnu.org/licenses/>.

// Range Governance system is the fundamental part of Range Core.

// NOTE: It's not allowed to change the compiler due to byte-to-byte
//       match requirement.
pragma solidity 0.5.16;
//pragma experimental SMTChecker;

import { GlobalConstants } from "./constants.sol";
import { IGovernedContract, GovernedContract } from "./GovernedContract.sol";
import { IBlockReward } from "./IBlockReward.sol";

/**
 * Genesis hardcoded version of SporkReward
 *
 * NOTE: it MUST NOT change after blockchain launch!
 */
contract StakerRewardV1 is
    GlobalConstants,
    GovernedContract,
    IBlockReward
{
    // IGovernedContract
    //---------------------------------
    // solium-disable-next-line no-empty-blocks
    constructor(address _proxy) public GovernedContract(_proxy) {}

    // IBlockReward
    //---------------------------------
    function reward() external payable {
        block.coinbase.transfer(msg.value);
    }

    function getReward(uint _blockNumber)
        external view
        returns(uint amount)
    {
        if (_blockNumber > 0) {
            amount = REWARD_STAKER_V1;
        }
    }

    // Safety
    //---------------------------------
    function () external payable {
        revert("Not supported");
    }
}
