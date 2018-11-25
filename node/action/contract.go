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
}

type Install struct {
	//ToDo: all the data you need to install contract
	Name    string
	Version version.Version
	Script  []byte
}

type Execute struct {
	//ToDo: all the data you need to execute contract
	Name    string
	Version version.Version
	Script  []byte
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
	if transaction.Function == INSTALL {
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
	}

	if transaction.Function == EXECUTE {
		if transaction.Owner == nil {
			log.Debug("Missing Data", "transaction owner", transaction.Owner)
			return status.MISSING_DATA
		}

		if transaction.Data == nil {
			log.Debug("Missing Data", "transaction data", transaction)
			return status.MISSING_DATA
		}

		executeData := transaction.Data.(Execute)
		if executeData.Name == "" {
			log.Debug("Missing Data", "name", executeData.Name)
			return status.MISSING_DATA
		}
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

func Convert(installData Install) (string, version.Version, data.Script) {
	name := installData.Name
	version := installData.Version
	script := data.Script{
		Script: installData.Script,
	}
	return name, version, script
}

func (transaction *Contract) ProcessDeliver(app interface{}) status.Code {
	log.Debug("Processing Smart Contract Transaction for DeliverTx")

	if transaction.Function == INSTALL {
		owner := transaction.Owner
		installData := transaction.Data.(Install)
		name, version, script := Convert(installData)

		smartContracts := GetSmartContracts(app)
		var scriptRecord *data.Scripts
		raw := smartContracts.Get(owner)
		if raw == nil {
			scriptRecord = data.NewScripts()
		} else {
			scriptRecord = raw.(*data.Scripts)
		}
		scriptRecord.Set(name, version, script)
		session := smartContracts.Begin()
		session.Set(owner, scriptRecord)
		session.Commit()
	}

	if transaction.Function == EXECUTE {
		owner := transaction.Owner
		executeData := transaction.Data.(Execute)
		smartContracts := GetSmartContracts(app)
		raw := smartContracts.Get(owner)
		if raw != nil {
			scriptRecord := raw.(*data.Scripts)
			versions := scriptRecord.Name[executeData.Name]
			script := versions.Version[executeData.Version.String()]
			RunScript(script.Script)
		}
	}

	return status.SUCCESS
}

func (transaction *Contract) Resolve(app interface{}) Commands {
	return []Command{}
}

func RunScript(script []byte) {
	log.Debug("Execute script", "script", string(script))
}
