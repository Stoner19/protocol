/*
	Copyright 2017-2018 OneLedger

	Generic Access to the different database in App, for cross-pkg access
*/

package app

import "github.com/Oneledger/protocol/node/log"

// Access to the local persistent databases
func (app Application) GetAdmin() interface{} {
	return app.Admin
}

// Access to the local persistent databases
func (app Application) GetStatus() interface{} {
	return app.Status
}

// Access to the local persistent databases
func (app Application) GetIdentities() interface{} {
	return app.Identities
}

// Access to the local persistent databases
func (app Application) GetAccounts() interface{} {
	return app.Accounts
}

// Access to the local persistent databases
func (app Application) GetBalances() interface{} {
	return app.Balances
}

func (app Application) GetChainID() interface{} {
	return ChainId
}

func (app Application) GetEvent() interface{} {
	return app.Event
}

func (app Application) GetContract() interface{} {
	return app.Contract
}

func (app Application) GetSmartContract() interface{} {
	return app.SmartContract
}

func (app Application) GetValidators() interface{} {
	log.Dump("Pat2", app.Validators)
	return app.Validators
}

func (app Application) GetSequence() interface{} {
	return app.Sequence
}
