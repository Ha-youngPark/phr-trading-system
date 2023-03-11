/*
 * SPDX-License-Identifier: Apache-2.0
 */

package phr

import ledgerapi "github.com/Ha-youngPark/phr-trading-system/phr/organization/hospital/contract-go/ledger-api"


// ListInterface defines functionality needed
// to interact with the world state on behalf
// of a phr
type ListInterface interface {
	AddPHR(*PHR) error
	GetPHR(string, string) (*PHR, error)
	UpdatePHR(*PHR) error
}

type list struct {
	stateList ledgerapi.StateListInterface
}

func (phrl *list) AddPHR(phr *PHR) error {
	return phrl.stateList.AddState(phr)
}

func (phrl *list) GetPHR(issuer string, phrNumber string) (*PHR, error) {
	phr := new(PHR)

	err := phrl.stateList.GetState(CreatePHRKey(issuer, phrNumber), phr)

	if err != nil {
		return nil, err
	}

	return phr, nil
}

func (phrl *list) UpdatePHR(phr *PHR) error {
	return phrl.stateList.UpdateState(phr)
}

// NewList create a new list from context
func newList(ctx TransactionContextInterface) *list {
	stateList := new(ledgerapi.StateList)
	stateList.Ctx = ctx
	stateList.Name = "org.phrnet.phrlist" 
	stateList.Deserialize = func(bytes []byte, state ledgerapi.StateInterface) error {
		return Deserialize(bytes, state.(*PHR))
	}

	list := new(list)
	list.stateList = stateList

	return list
}
