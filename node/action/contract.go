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
	Version string
	Name    string
}

type Execute struct {
	//ToDo: all the data you need to execute contract
}

func init() {
	serial.Register(Contract{})
}

func (transaction *Contract) TransactionType() Type {
	return transaction.Base.Type
}

func (transaction *Contract) Validate() status.Code {
	log.Debug("Validating Contract Transaction")

	if transaction.Data == nil {
		log.Debug("Missing Data", "contract", transaction)
		return status.MISSING_DATA
	}

	//if transaction.Function == nil {
	//	log.Debug("Missing Gas", "contract", transaction)
	//	return status.MISSING_DATA
	//}

	return status.SUCCESS
}

func (transaction *Contract) ProcessCheck(app interface{}) status.Code {
	log.Debug("Processing Contract Transaction for CheckTx")

	//if !CheckAmounts(app, transaction.Data, transaction.Function) {
	//	log.Debug("FAILED", "inputs", transaction.Data, "outputs", transaction.Function)
	//	return status.INVALID
	//return status.SUCCESS
	//}

	// TODO: Validate the transaction against the UTXO database, check tree
	balances := GetBalances(app)
	_ = balances

	return status.SUCCESS
}

func (transaction *Contract) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *Contract) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Contract Transaction for DeliverTx")

	//if !CheckAmounts(app, transaction.Data, transaction.Function) {
	//	return status.INVALID
	//}

	//	balances := GetBalances(app)

	// Update the database to the final set of entries
	//for _, entry := range transaction.Outputs {
	//	var balance *data.Balance
	//	result := balances.Get(entry.AccountKey)
	//	if result == nil {
	//		tmp := data.NewBalance()
	//		result = &tmp
	//	}
	//	balance = result
	//	balance.SetAmmount(entry.Amount)
	//
	//		balances.Set(entry.AccountKey, *balance)
	//	}

	return status.SUCCESS
}

func (transaction *Contract) Resolve(app interface{}) Commands {
	return []Command{}
}
