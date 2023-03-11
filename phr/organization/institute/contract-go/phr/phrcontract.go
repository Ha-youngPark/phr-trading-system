/*
 * SPDX-License-Identifier: Apache-2.0
 */

package phr

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Contract chaincode that defines
// the business logic for managing phr
type Contract struct {
	contractapi.Contract
}

// Instantiate does nothing
func (c *Contract) Instantiate() {
	fmt.Println("Instantiated")
}

// Issue creates a new phr and stores it in the world state
func (c *Contract) Issue(ctx TransactionContextInterface, issuer string, phrNumber string, issueDateTime string, maturityDateTime string, faceValue int) (*PHR, error) {
	phr := PHR{PHRNumber: phrNumber, Issuer: issuer, IssueDateTime: issueDateTime, FaceValue: faceValue, MaturityDateTime: maturityDateTime, Owner: issuer}
	phr.SetIssued()

	err := ctx.GetPHRList().AddPHR(&phr)

	if err != nil {
		return nil, err
	}

	return &phr, nil
}

// Buy updates a phr to be in trading status and sets the new owner
func (c *Contract) Buy(ctx TransactionContextInterface, issuer string, phrNumber string, currentOwner string, newOwner string, price int, purchaseDateTime string) (*PHR, error) {
	phr, err := ctx.GetPHRList().GetPHR(issuer, phrNumber)

	if err != nil {
		return nil, err
	}

	if phr.Owner != currentOwner {
		return nil, fmt.Errorf("PHR %s:%s is not owned by %s", issuer, phrNumber, currentOwner)
	}

	if phr.IsIssued() {
		phr.SetTrading()
	}

	if !phr.IsTrading() {
		return nil, fmt.Errorf("PHR %s:%s is not trading. Current state = %s", issuer, phrNumber, phr.GetState())
	}

	phr.Owner = newOwner

	err = ctx.GetPHRList().UpdatePHR(phr)

	if err != nil {
		return nil, err
	}

	return phr, nil
}

// Redeem updates a phr status to be expired
func (c *Contract) Expire(ctx TransactionContextInterface, issuer string, phrNumber string, expiringOwner string, expireDateTime string) (*PHR, error) {
	phr, err := ctx.GetPHRList().GetPHR(issuer, phrNumber)

	if err != nil {
		return nil, err
	}

	if phr.Owner != expiringOwner {
		return nil, fmt.Errorf("PHR %s:%s is not owned by %s", issuer, phrNumber, expiringOwner)
	}

	if phr.IsExpired() {
		return nil, fmt.Errorf("PHR %s:%s is already expired", issuer, phrNumber)
	}

	phr.Owner = phr.Issuer
	phr.SetExpired()

	err = ctx.GetPHRList().UpdatePHR(phr)

	if err != nil {
		return nil, err
	}

	return phr, nil
}
