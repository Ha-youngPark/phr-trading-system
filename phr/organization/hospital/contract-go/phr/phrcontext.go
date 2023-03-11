/*
 * SPDX-License-Identifier: Apache-2.0
 */

package phr

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// TransactionContextInterface an interface to
// describe the minimum required functions for
// a transaction context in the phr
type TransactionContextInterface interface {
	contractapi.TransactionContextInterface
	GetPHRList() ListInterface
}

// TransactionContext implementation of
// TransactionContextInterface for use with
// phr contract
type TransactionContext struct {
	contractapi.TransactionContext
	phrList *list
}

// GetPHRList return phr list
func (tc *TransactionContext) GetPHRList() ListInterface {
	if tc.phrList == nil {
		tc.phrList = newList(tc)
	}

	return tc.phrList
}
