// Copyright 2018 The Range Core Authors
// Copyright 2014 The go-ethereum Authors
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

package tests

import (
	"os"
	"testing"

	"range/core/gen3/core/vm"
)

func TestVM(t *testing.T) {
	if val, ok := os.LookupEnv("SKIP_KNOWN_FAIL"); ok && val == "1" {
		t.Skip("unit test is broken: conditional test skipping activated")
	}
	t.Parallel()
	vmt := new(testMatcher)
	vmt.slow("^vmPerformance")
	vmt.fails("^vmSystemOperationsTest.json/createNameRegistrator$", "fails without parallel execution")

	vmt.walk(t, vmTestDir, func(t *testing.T, name string, test *VMTest) {
		withTrace(t, test.json.Exec.GasLimit, func(vmconfig vm.Config) error {
			return vmt.checkFailure(t, name, test.Run(vmconfig))
		})
	})
}
