/*
 * SPDX-License-Identifier: Apache-2.0
 */

package phr

import (
	"errors"
	"testing"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// #########
// HELPERS
// #########
type MockPHRList struct {
	mock.Mock
}

func (mpl *MockPHRList) AddPHR(phr *PHR) error {
	args := mpl.Called(phr)

	return args.Error(0)
}

func (mpl *MockPHRList) GetPHR(issuer string, phrnumber string) (*PHR, error) {
	args := mpl.Called(issuer, phrnumber)

	return args.Get(0).(*PHR), args.Error(1)
}

func (mpl *MockPHRList) UpdatePHR(phr *PHR) error {
	args := mpl.Called(phr)

	return args.Error(0)
}

type MockTransactionContext struct {
	contractapi.TransactionContext
	phrList *MockPHRList
}

func (mtc *MockTransactionContext) GetPHRList() ListInterface {
	return mtc.phrList
}

func resetPHR(phr *PHR) {
	phr.Owner = "someowner"
	phr.SetTrading()
}

// #########
// TESTS
// #########

func TestIssue(t *testing.T) {
	var phr *PHR
	var err error

	mpl := new(MockPHRList)
	ctx := new(MockTransactionContext)
	ctx.phrList = mpl

	contract := new(Contract)

	var sentPHR *PHR

	mpl.On("AddPHR", mock.MatchedBy(func(phr *PHR) bool { sentPHR = phr; return phr.Issuer == "someissuer" })).Return(nil)
	mpl.On("AddPHR", mock.MatchedBy(func(phr *PHR) bool { sentPHR = phr; return phr.Issuer == "someotherissuer" })).Return(errors.New("AddPHR error"))

	expectedPHR := PHR{PHRNumber: "somephr", Issuer: "someissuer", IssueDateTime: "someissuedate", FaceValue: 1000, MaturityDateTime: "somematuritydate", Owner: "someissuer", state: 1}
	phr, err = contract.Issue(ctx, "someissuer", "somephr", "someissuedate", "somematuritydate", 1000)
	assert.Nil(t, err, "should not error when add phr does not error")
	assert.Equal(t, sentPHR, phr, "should send the same phr as it returns to add phr")
	assert.Equal(t, expectedPHR, *phr, "should correctly configure phr")

	phr, err = contract.Issue(ctx, "someotherissuer", "somephr", "someissuedate", "somematuritydate", 1000)
	assert.EqualError(t, err, "AddPHR error", "should return error when add phr fails")
	assert.Nil(t, phr, "should not return phr when fails")
}

func TestBuy(t *testing.T) {
	var phr *PHR
	var err error

	mpl := new(MockPHRList)
	ctx := new(MockTransactionContext)
	ctx.phrList = mpl

	contract := new(Contract)

	wsPHR := new(PHR)
	resetPHR(wsPHR)

	var sentPHR *PHR
	var emptyPHR *PHR
	shouldError := false

	mpl.On("GetPHR", "someissuer", "somephr").Return(wsPHR, nil)
	mpl.On("GetPHR", "someotherissuer", "someotherphr").Return(emptyPHR, errors.New("GetPHR error"))
	mpl.On("UpdatePHR", mock.MatchedBy(func(phr *PHR) bool { return shouldError })).Return(errors.New("UpdatePHR error"))
	mpl.On("UpdatePHR", mock.MatchedBy(func(phr *PHR) bool { sentPHR = phr; return !shouldError })).Return(nil)

	phr, err = contract.Buy(ctx, "someotherissuer", "someotherphr", "someowner", "someotherowner", 100, "2019-12-10:10:00")
	assert.EqualError(t, err, "GetPHR error", "should return error when GetPHR errors")
	assert.Nil(t, phr, "should return nil for phr when GetPHR errors")

	phr, err = contract.Buy(ctx, "someissuer", "somephr", "someotherowner", "someowner", 100, "2019-12-10:10:00")
	assert.EqualError(t, err, "PHR someissuer:somephr is not owned by someotherowner", "should error when sent owner not correct")
	assert.Nil(t, phr, "should not return phr for bad owner error")

	resetPHR(wsPHR)
	wsPHR.SetExpired()
	phr, err = contract.Buy(ctx, "someissuer", "somephr", "someowner", "someotherowner", 100, "2019-12-10:10:00")
	assert.EqualError(t, err, "PHR someissuer:somephr is not trading. Current state = EXPIRED")
	assert.Nil(t, phr, "should not return phr for bad state error")

	resetPHR(wsPHR)
	shouldError = true
	phr, err = contract.Buy(ctx, "someissuer", "somephr", "someowner", "someotherowner", 100, "2019-12-10:10:00")
	assert.EqualError(t, err, "UpdatePHR error", "should error when update phr fails")
	assert.Nil(t, phr, "should not return phr for bad state error")
	shouldError = false

	resetPHR(wsPHR)
	wsPHR.SetIssued()
	phr, err = contract.Buy(ctx, "someissuer", "somephr", "someowner", "someotherowner", 100, "2019-12-10:10:00")
	assert.Nil(t, err, "should not error when good phr and owner")
	assert.Equal(t, "someotherowner", phr.Owner, "should update the owner of the phr")
	assert.True(t, phr.IsTrading(), "should mark issued phr as trading")
	assert.Equal(t, sentPHR, phr, "should update same phr as it returns in the world state")
}

func TestExpire(t *testing.T) {
	var phr *PHR
	var err error

	mpl := new(MockPHRList)
	ctx := new(MockTransactionContext)
	ctx.phrList = mpl

	contract := new(Contract)

	var sentPHR *PHR
	wsPHR := new(PHR)
	resetPHR(wsPHR)

	var emptyPHR *PHR
	shouldError := false

	mpl.On("GetPHR", "someissuer", "somephr").Return(wsPHR, nil)
	mpl.On("GetPHR", "someotherissuer", "someotherphr").Return(emptyPHR, errors.New("GetPHR error"))
	mpl.On("UpdatePHR", mock.MatchedBy(func(phr *PHR) bool { return shouldError })).Return(errors.New("UpdatePHR error"))
	mpl.On("UpdatePHR", mock.MatchedBy(func(phr *PHR) bool { sentPHR = phr; return !shouldError })).Return(nil)

	phr, err = contract.Expire(ctx, "someotherissuer", "someotherphr", "someowner", "2021-12-10:10:00")
	assert.EqualError(t, err, "GetPHR error", "should error when GetPHR errors")
	assert.Nil(t, phr, "should not return phr when GetPHR errors")

	phr, err = contract.Expire(ctx, "someissuer", "somephr", "someotherowner", "2021-12-10:10:00")
	assert.EqualError(t, err, "PHR someissuer:somephr is not owned by someotherowner", "should error when phr owned by someone else")
	assert.Nil(t, phr, "should not return phr when errors as owned by someone else")

	resetPHR(wsPHR)
	wsPHR.SetExpired()
	phr, err = contract.Expire(ctx, "someissuer", "somephr", "someowner", "2021-12-10:10:00")
	assert.EqualError(t, err, "PHR someissuer:somephr is already expired", "should error when phr already expired")
	assert.Nil(t, phr, "should not return phr when errors as already expired")

	shouldError = true
	resetPHR(wsPHR)
	phr, err = contract.Expire(ctx, "someissuer", "somephr", "someowner", "2021-12-10:10:00")
	assert.EqualError(t, err, "UpdatePHR error", "should error when update phr errors")
	assert.Nil(t, phr, "should not return phr when UpdatePHR errors")
	shouldError = false

	resetPHR(wsPHR)
	phr, err = contract.Expire(ctx, "someissuer", "somephr", "someowner", "2021-12-10:10:00")
	assert.Nil(t, err, "should not error on good expired")
	assert.True(t, phr.IsExpired(), "should return expired phr")
	assert.Equal(t, sentPHR, phr, "should update same phr as it returns in the world state")
}
