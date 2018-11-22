/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"github.com/Oneledger/protocol/node/version"
)

type ContractFunction int

const (
	INSTALL ContractFunction = iota
	EXECUTE
)

// Synchronize a swap between two users
type Contract struct {
	Base
	Data     ContractData     `json:"data"`
	Function ContractFunction `json:"function"`
	Gas      data.Coin        `json:"gas"`
	Fee      data.Coin        `json:"fee"`
}

type ContractData interface {
	//validate() status.Code
	//resolve(interface{}, swapStageType) Commands
}

type Install struct {
	//ToDo: all the data you need to install contract
	Name    string
	Version version.Version
	Script  []byte
}

type Execute struct {
	//ToDo: all the data you need to execute contract
}

func init() {
	serial.Register(Contract{})
	serial.Register(Install{})
	serial.Register(Execute{})

	var prototype ContractData
	serial.RegisterInterface(&prototype)
}

func (transaction *Contract) TransactionType() Type {
	return transaction.Base.Type
}

func (transaction *Contract) Validate() status.Code {
	log.Debug("Validating Smart Contract Transaction")

	//check that the data supplied is valid and no security problems

	if transaction.Owner == nil {
		log.Debug("Missing Data", "transaction owner", transaction.Owner)
		return status.MISSING_DATA
	}

	if transaction.Data == nil {
		log.Debug("Missing Data", "transaction data", transaction)
		return status.MISSING_DATA
	}

	installData := transaction.Data.(Install)
	if installData.Name == "" {
		log.Debug("Missing Data", "name", installData.Name)
		return status.MISSING_DATA
	}

	if installData.Script == nil {
		log.Debug("Missing Data", "script", installData.Script)
		return status.MISSING_DATA
	}

	return status.SUCCESS
}

func (transaction *Contract) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Smart Contract Transaction for CheckTx")

	return status.SUCCESS
}

func (transaction *Contract) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *Contract) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Smart Contract Transaction for DeliverTx")

	owner := transaction.Owner
	installData := transaction.Data.(Install)

	smartContracts := GetSmartContracts(app)
	session := smartContracts.Begin()
	session.Set(owner, installData)
	session.Commit()

	return status.SUCCESS
}

func (transaction *Contract) Resolve(app interface{}) Commands {
	return []Command{}
}
