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

import { MasternodeTokenV2 } from "./MasternodeTokenV2.sol";
import { IGovernedProxy } from "./IGovernedProxy.sol";

contract MockMasternodeTokenV2
    is MasternodeTokenV2
{
    constructor(address _proxy, IGovernedProxy _registry_proxy)
        public
        MasternodeTokenV2(_proxy, _registry_proxy)
    // solium-disable-next-line no-empty-blocks
    {}

    function setBalance(address tokenOwner, uint balance) external {
        v1storage.setBalance(tokenOwner, balance, block.number);
    }
}
