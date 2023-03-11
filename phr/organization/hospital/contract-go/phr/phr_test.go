/*
 * SPDX-License-Identifier: Apache-2.0
 */

package phr

import (
	"testing"

	ledgerapi "github.com/Ha-youngPark/phr-trading-system/phr/organization/hospital/contract-go/ledger-api"
	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	assert.Equal(t, "ISSUED", ISSUED.String(), "should return string for issued")
	assert.Equal(t, "TRADING", TRADING.String(), "should return string for issued")
	assert.Equal(t, "EXPIRED", EXPIRED.String(), "should return string for issued")
	assert.Equal(t, "UNKNOWN", State(EXPIRED+1).String(), "should return unknown when not one of constants")
}

func TestCreatePHRKey(t *testing.T) {
	assert.Equal(t, ledgerapi.MakeKey("someissuer", "somephr"), CreatePHRKey("someissuer", "somephr"), "should return key comprised of passed values")
}

func TestGetState(t *testing.T) {
	phr := new(PHR)
	phr.state = ISSUED

	assert.Equal(t, ISSUED, phr.GetState(), "should return set state")
}

func TestSetIssued(t *testing.T) {
	phr := new(PHR)
	phr.SetIssued()
	assert.Equal(t, ISSUED, phr.state, "should set state to trading")
}

func TestSetTrading(t *testing.T) {
	phr := new(PHR)
	phr.SetTrading()
	assert.Equal(t, TRADING, phr.state, "should set state to trading")
}

func TestSetExpired(t *testing.T) {
	phr := new(PHR)
	phr.SetExpired()
	assert.Equal(t, EXPIRED, phr.state, "should set state to trading")
}

func TestIsIssued(t *testing.T) {
	phr := new(PHR)

	phr.SetIssued()
	assert.True(t, phr.IsIssued(), "should be true when status set to issued")

	phr.SetTrading()
	assert.False(t, phr.IsIssued(), "should be false when status not set to issued")
}

func TestIsTrading(t *testing.T) {
	phr := new(PHR)

	phr.SetTrading()
	assert.True(t, phr.IsTrading(), "should be true when status set to trading")

	phr.SetExpired()
	assert.False(t, phr.IsTrading(), "should be false when status not set to trading")
}

func TestIsExpired(t *testing.T) {
	phr := new(PHR)

	phr.SetExpired()
	assert.True(t, phr.IsExpired(), "should be true when status set to expired")

	phr.SetIssued()
	assert.False(t, phr.IsExpired(), "should be false when status not set to expired")
}

func TestGetSplitKey(t *testing.T) {
	phr := new(PHR)
	phr.PHRNumber = "somephr"
	phr.Issuer = "someissuer"

	assert.Equal(t, []string{"someissuer", "somephr"}, phr.GetSplitKey(), "should return issuer and phr number as split key")
}

func TestSerialize(t *testing.T) {
	phr := new(PHR)
	phr.PHRNumber = "somephr"
	phr.Issuer = "someissuer"
	phr.IssueDateTime = "sometime"
	phr.FaceValue = 1000
	phr.MaturityDateTime = "somelatertime"
	phr.Owner = "someowner"
	phr.state = TRADING

	bytes, err := phr.Serialize()
	assert.Nil(t, err, "should not error on serialize")
	assert.Equal(t, `{"phrNumber":"somephr","issuer":"someissuer","issueDateTime":"sometime","faceValue":1000,"maturityDateTime":"somelatertime","owner":"someowner","currentState":2,"class":"org.phrnet.phrlist","key":"someissuer:somephr"}`, string(bytes), "should return JSON formatted value")
}

func TestDeserialize(t *testing.T) {
	var phr *PHR
	var err error

	goodJSON := `{"phrNumber":"somephr","issuer":"someissuer","issueDateTime":"sometime","faceValue":1000,"maturityDateTime":"somelatertime","owner":"someowner","currentState":2,"class":"org.phrnet.phr","key":"someissuer:somephr"}`
	expectedPHR := new(PHR)
	expectedPHR.PHRNumber = "somephr"
	expectedPHR.Issuer = "someissuer"
	expectedPHR.IssueDateTime = "sometime"
	expectedPHR.FaceValue = 1000
	expectedPHR.MaturityDateTime = "somelatertime"
	expectedPHR.Owner = "someowner"
	expectedPHR.state = TRADING
	phr = new(PHR)
	err = Deserialize([]byte(goodJSON), phr)
	assert.Nil(t, err, "should not return error for deserialize")
	assert.Equal(t, expectedPHR, phr, "should create expected phr")

	badJSON := `{"phrNumber":"somephr","issuer":"someissuer","issueDateTime":"sometime","faceValue":"NaN","maturityDateTime":"somelatertime","owner":"someowner","currentState":2,"class":"org.phrnet.phrlist","key":"someissuer:somephr"}`
	phr = new(PHR)
	err = Deserialize([]byte(badJSON), phr)
	assert.EqualError(t, err, "Error deserializing phr. json: cannot unmarshal string into Go struct field jsonPHR.faceValue of type int", "should return error for bad data")
}
