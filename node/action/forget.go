/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/log"
)

type Forget struct {
	Base

	Target string `json:"target"`
}

func (transaction *Forget) Validate() err.Code {
	log.Debug("Validating Forget Transaction")
	return err.SUCCESS
}

func (transaction *Forget) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Forget Transaction for CheckTx")
	return err.SUCCESS
}

func (transaction *Forget) ShouldProcess(app interface{}) bool {
	return true
}

func (transaction *Forget) ProcessDeliver(app interface{}) err.Code {

	commands := transaction.Resolve(app)

	//before loop of execute, lastResult is nil
	return commands.Execute(app)
}

func (transaction *Forget) Resolve(app interface{}) Commands {
	return []Command{}
}

