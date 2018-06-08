/*
	Copyright 2017 - 2018 OneLedger

	Register this identity with the other nodes. As an externl identity
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

// Register an identity with the chain
type Register struct {
	Base

	Identity string
}

func (transaction Register) Validate() err.Code {
	log.Debug("Validating Register Transaction")

	// TODO: Make sure all of the parameters are there
	// TODO: Check all signatures and keys
	// TODO: Vet that the sender has the values
	return err.SUCCESS
}

// Test to see if the identity already exists
func (transaction Register) ProcessCheck(app interface{}) err.Code {
	log.Debug("Processing Register Transaction for CheckTx")

	identities := GetIdentities(app)
	id, errs := identities.FindName(transaction.Identity)

	if errs != err.SUCCESS {
		return errs
	}

	if id == nil {
		log.Debug("Success, can add new Identity", "id", id)
		return err.SUCCESS
	}

	log.Debug("Identity already exists", "id", id)

	// TODO: Not necessarily a failure, since this identity might be local
	return err.SUCCESS
}

// Add the identity into the database as external, don't overwrite a local identity
func (transaction Register) ProcessDeliver(app interface{}) err.Code {
	log.Debug("Processing Register Transaction for DeliverTx")

	identities := GetIdentities(app)
	entry, errs := identities.FindName(transaction.Identity)

	if errs != err.SUCCESS {
		return errs
	}

	if entry != nil {
		/*
			if !entry.IsExternal() {
				return err.SUCCESS
			}
		*/
		log.Debug("Ignoring Duplicate Identity")
	} else {
		identities.Add(id.NewIdentity(transaction.Identity, "Contact Information", true))
	}

	log.Info("Updated External Identity Reference!!!", "id", transaction.Identity)

	return err.SUCCESS
}

// Given a transaction, expand it into a list of Commands to execute against various chains.
func (transaction Register) Expand(app interface{}) Commands {
	// TODO: Table-driven mechanics, probably elsewhere
	return []Command{}
}
