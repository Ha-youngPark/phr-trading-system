/*
 * SPDX-License-Identifier: Apache-2.0
 */

package phr

import (
	"errors"
	"testing"

	ledgerapi "github.com/Ha-youngPark/phr-trading-system/phr/organization/hospital/contract-go/ledger-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #########
// HELPERS
// #########

type MockStateList struct {
	mock.Mock
}

func (msl *MockStateList) AddState(state ledgerapi.StateInterface) error {
	args := msl.Called(state)

	return args.Error(0)
}

func (msl *MockStateList) GetState(key string, state ledgerapi.StateInterface) error {
	args := msl.Called(key, state)

	state.(*PHR).PHRNumber = "somephr"

	return args.Error(0)
}

func (msl *MockStateList) UpdateState(state ledgerapi.StateInterface) error {
	args := msl.Called(state)

	return args.Error(0)
}

// #########
// TESTS
// #########

func TestAddPHR(t *testing.T) {
	phr := new(PHR)

	list := new(list)
	msl := new(MockStateList)
	msl.On("AddState", phr).Return(errors.New("Called add state correctly"))
	list.stateList = msl

	err := list.AddPHR(phr)
	assert.EqualError(t, err, "Called add state correctly", "should call state list add state with phr")
}

func TestGetPHR(t *testing.T) {
	var phr *PHR
	var err error

	list := new(list)
	msl := new(MockStateList)
	msl.On("GetState", CreatePHRKey("someissuer", "somephr"), mock.MatchedBy(func(state ledgerapi.StateInterface) bool { _, ok := state.(*PHR); return ok })).Return(nil)
	msl.On("GetState", CreatePHRKey("someotherissuer", "someotherphr"), mock.MatchedBy(func(state ledgerapi.StateInterface) bool { _, ok := state.(*PHR); return ok })).Return(errors.New("GetState error"))
	list.stateList = msl

	phr, err = list.GetPHR("someissuer", "somephr")
	assert.Nil(t, err, "should not error when get state on state list does not error")
	assert.Equal(t, phr.PHRNumber, "somephr", "should use state list GetState to fill phr")

	phr, err = list.GetPHR("someotherissuer", "someotherphr")
	assert.EqualError(t, err, "GetState error", "should return error when state list get state errors")
	assert.Nil(t, phr, "should not return phr on error")
}

func TestUpdatePHR(t *testing.T) {
	phr := new(PHR)

	list := new(list)
	msl := new(MockStateList)
	msl.On("UpdateState", phr).Return(errors.New("Called update state correctly"))
	list.stateList = msl

	err := list.UpdatePHR(phr)
	assert.EqualError(t, err, "Called update state correctly", "should call state list update state with phr")
}

func TestNewStateList(t *testing.T) {
	ctx := new(TransactionContext)
	list := newList(ctx)
	stateList, ok := list.stateList.(*ledgerapi.StateList)

	assert.True(t, ok, "should make statelist of type ledgerapi.StateList")
	assert.Equal(t, ctx, stateList.Ctx, "should set the context to passed context")
	assert.Equal(t, "org.phrnet.phrlist", stateList.Name, "should set the name for the list")

	expectedErr := Deserialize([]byte("bad json"), new(PHR))
	err := stateList.Deserialize([]byte("bad json"), new(PHR))
	assert.EqualError(t, err, expectedErr.Error(), "should call Deserialize when stateList.Deserialize called")
}
