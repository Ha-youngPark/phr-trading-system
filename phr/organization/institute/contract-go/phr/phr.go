/*
 * SPDX-License-Identifier: Apache-2.0
 */

package phr

import (
	"encoding/json"
	"fmt"

	ledgerapi "github.com/Ha-youngPark/phr-trading-system/phr/organization/institute/contract-go/ledger-api"
)

// State enum for phr state property
type State uint

const (
	// ISSUED state for when a phr has been issued
	ISSUED State = iota + 1
	// TRADING state for when a phr is trading
	TRADING
	// EXPIRED state for when a phr has been expired
	EXPIRED
)

func (state State) String() string {
	names := []string{"ISSUED", "TRADING", "EXPIRED"}

	if state < ISSUED || state > EXPIRED {
		return "UNKNOWN"
	}

	return names[state-1]
}

// CreatePHRKey creates a key for phrs
func CreatePHRKey(issuer string, phrNumber string) string {
	return ledgerapi.MakeKey(issuer, phrNumber)
}

// Used for managing the fact status is private but want it in world state
type PHRAlias PHR
type jsonPHR struct {
	*PHRAlias
	State State  `json:"currentState"`
	Class string `json:"class"`
	Key   string `json:"key"`
}

// PHR defines a phr
type PHR struct {
	PHRNumber      string `json:"phrNumber"`
	Issuer           string `json:"issuer"`
	IssueDateTime    string `json:"issueDateTime"`
	FaceValue        int    `json:"faceValue"`
	MaturityDateTime string `json:"maturityDateTime"`
	Owner            string `json:"owner"`
	state            State  `metadata:"currentState"`
	class            string `metadata:"class"`
	key              string `metadata:"key"`
}

// UnmarshalJSON special handler for managing JSON marshalling
func (phr *PHR) UnmarshalJSON(data []byte) error {
	jphr := jsonPHR{PHRAlias: (*PHRAlias)(phr)}

	err := json.Unmarshal(data, &jphr)

	if err != nil {
		return err
	}

	phr.state = jphr.State

	return nil
}

// MarshalJSON special handler for managing JSON marshalling
func (phr PHR) MarshalJSON() ([]byte, error) {
	jphr := jsonPHR{PHRAlias: (*PHRAlias)(&phr), State: phr.state, Class: "org.phrnet.phrlist", Key: ledgerapi.MakeKey(phr.Issuer, phr.PHRNumber)}

	return json.Marshal(&jphr)
}

// GetState returns the state
func (phr *PHR) GetState() State {
	return phr.state
}

// SetIssued returns the state to issued
func (phr *PHR) SetIssued() {
	phr.state = ISSUED
}

// SetTrading sets the state to trading
func (phr *PHR) SetTrading() {
	phr.state = TRADING
}

// SetRedeemed sets the state to redeemed
func (phr *PHR) SetExpired() {
	phr.state = EXPIRED
}

// IsIssued returns true if state is issued
func (phr *PHR) IsIssued() bool {
	return phr.state == ISSUED
}

// IsTrading returns true if state is trading
func (phr *PHR) IsTrading() bool {
	return phr.state == TRADING
}

// IsRedeemed returns true if state is redeemed
func (phr *PHR) IsExpired() bool {
	return phr.state == EXPIRED
}

// GetSplitKey returns values which should be used to form key
func (phr *PHR) GetSplitKey() []string {
	return []string{phr.Issuer, phr.PHRNumber}
}

// Serialize formats the commercial paper as JSON bytes
func (phr *PHR) Serialize() ([]byte, error) {
	return json.Marshal(phr)
}

// Deserialize formats the commercial paper from JSON bytes
func Deserialize(bytes []byte, phr *PHR) error {
	err := json.Unmarshal(bytes, phr)

	if err != nil {
		return fmt.Errorf("Error deserializing phr. %s", err.Error())
	}

	return nil
}
