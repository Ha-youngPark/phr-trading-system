/*
 * SPDX-License-Identifier: Apache-2.0
 */

package phr

import (
	"testing"

	ledgerapi "github.com/Ha-youngPark/phr-trading-system/phr/organization/hospital/contract-go/ledger-api"
	"github.com/stretchr/testify/assert"
)

func TestGetPHRList(t *testing.T) {
	var tc *TransactionContext
	var expectedPHRList *list

	tc = new(TransactionContext)
	expectedPHRList = newList(tc)
	actualList := tc.GetPHRList().(*list)
	assert.Equal(t, expectedPHRList.stateList.(*ledgerapi.StateList).Name, actualList.stateList.(*ledgerapi.StateList).Name, "should configure phr list when one not already configured")

	tc = new(TransactionContext)
	expectedPHRList = new(list)
	expectedStateList := new(ledgerapi.StateList)
	expectedStateList.Ctx = tc
	expectedStateList.Name = "existing phr list"
	expectedPHRList.stateList = expectedStateList
	tc.phrList = expectedPHRList
	assert.Equal(t, expectedPHRList, tc.GetPHRList(), "should return set phr list when already set")
}
