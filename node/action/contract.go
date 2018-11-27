/*
	Copyright 2017-2018 OneLedger

	Handle Smart Contarct Install and Execute functions
*/
package action

import (
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/data"
	"github.com/Oneledger/protocol/node/global"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
	"github.com/Oneledger/protocol/node/version"
)

type ContractFunction int

const (
	INSTALL ContractFunction = iota
	EXECUTE
	COMPARE
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
}

type Compare struct {
	//ToDo: all the data you need to execute contract
	Name    string
	Version version.Version
	Results string
}

func init() {
	serial.Register(Contract{})
	serial.Register(Install{})
	serial.Register(Execute{})
	serial.Register(Compare{})

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
			log.Debug("Smart Contract Missing Data", "transaction owner", transaction.Owner)
			return status.MISSING_DATA
		}

		if transaction.Data == nil {
			log.Debug("Smart Contract Missing Data", "transaction data", transaction)
			return status.MISSING_DATA
		}

		installData := transaction.Data.(Install)
		if installData.Name == "" {
			log.Debug("Smart Contract Missing Data", "name", installData.Name)
			return status.MISSING_DATA
		}

		if installData.Version.String() == "" {
			log.Debug("Smart Contract Missing Data", "version", installData.Version)
			return status.MISSING_DATA
		}

		if installData.Script == nil {
			log.Debug("Smart Contract Missing Data", "script", installData.Script)
			return status.MISSING_DATA
		}
	}

	if transaction.Function == EXECUTE {
		if transaction.Owner == nil {
			log.Debug("Smart Contract Missing Data", "transaction owner", transaction.Owner)
			return status.MISSING_DATA
		}

		if transaction.Data == nil {
			log.Debug("Smart Contract Missing Data", "transaction data", transaction)
			return status.MISSING_DATA
		}

		executeData := transaction.Data.(Execute)
		if executeData.Name == "" {
			log.Debug("Smart Contract Missing Data", "name", executeData.Name)
			return status.MISSING_DATA
		}

		if executeData.Version.String() == "" {
			log.Debug("Smart Contract Missing Data", "version", executeData.Version)
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

	//var result Transaction
	switch transaction.Function {
	case INSTALL:
		transaction.Install(app)
	case EXECUTE:
		transaction.Execute(app)
	case COMPARE:
		transaction.Compare(app)
	default:
		return status.INVALID
	}

	//if result != nil {
	//	log.Debug("JustBeforeBroadcastTransaction", "result", result)
	//	BroadcastTransaction(SMART_CONTRACT, result, false)
	//}
	return status.SUCCESS
}

func (transaction *Contract) Resolve(app interface{}) Commands {
	return []Command{}
}

func (transaction *Contract) Install(app interface{}) {
	owner := transaction.Owner
	installData := transaction.Data.(Install)
	name, version, script := Convert(installData)

	smartContracts := GetSmartContracts(app)
	var scriptRecords *data.ScriptRecords
	raw := smartContracts.Get(owner)
	if raw == nil {
		scriptRecords = data.NewScriptRecords()
	} else {
		scriptRecords = raw.(*data.ScriptRecords)
	}
	scriptRecords.Set(name, version, script)
	session := smartContracts.Begin()
	session.Set(owner, scriptRecords)
	session.Commit()
}

func (transaction *Contract) Execute(app interface{}) Transaction {
	//figure out the it node and run execute once
	//needs to take the result and put it in a new contract transaction and broadcast async

	validatorList := id.GetValidators(app)
	log.Dump("Pat3", validatorList)
	log.Debug("Execute NodeName", "validatorList", validatorList)
	selectedValidatorIdentity := validatorList.SelectedValidator
	log.Debug("Execute NodeName", "nodeName", selectedValidatorIdentity.NodeName)
	log.Debug("Execute NodeName", "globalNodeName", global.Current.NodeName)
	if global.Current.NodeName == selectedValidatorIdentity.NodeName {
		executeData := transaction.Data.(Execute)
		smartContracts := GetSmartContracts(app)
		raw := smartContracts.Get(transaction.Owner)
		if raw != nil {
			scriptRecords := raw.(*data.ScriptRecords)
			versions := scriptRecords.Name[executeData.Name]
			script := versions.Version[executeData.Version.String()]
			resultRunScript := RunScript(script.Script)
			if resultRunScript != "" {
				resultCompare := transaction.CreateCompareRequest(app, executeData.Name, executeData.Version, resultRunScript)
				if resultCompare != nil {
					//TODO: check this later
					comm.Broadcast(resultCompare)
					//BroadcastTransaction(SMARTCONTRACT, resultCompare, false)
					//return resultCompare
				}
			}

		}
	}
	return nil
}

func (transaction *Contract) Compare(app interface{}) status.Code {
	//see that new transaction, run the engine, if the results are the same, return success, otherwise fail
	log.Debug("InsideCompare", "transaction", transaction)
	compareData := transaction.Data.(Compare)
	smartContracts := GetSmartContracts(app)
	raw := smartContracts.Get(transaction.Owner)
	if raw != nil {
		scriptRecords := raw.(*data.ScriptRecords)
		versions := scriptRecords.Name[compareData.Name]
		script := versions.Version[compareData.Version.String()]
		resultRunScript := RunScript(script.Script)
		if resultRunScript == compareData.Results {
			return status.SUCCESS
		}
	}
	return status.INVALID
}

func RunScript(script []byte) string {
	log.Debug("Smart Contract Execute script", "script", string(script))
	return "Ta-dah"
}

func (transaction *Contract) CreateCompareRequest(app interface{}, name string, version version.Version, resultRunScript string) []byte {

	chainId := GetChainID(app)

	fee := data.NewCoin(0, "OLT")
	gas := data.NewCoin(0, "OLT")

	next := id.NextSequence(app, transaction.Owner)

	inputs := Compare{
		Name:    name,
		Version: version,
		Results: resultRunScript,
	}

	// Create base transaction
	compare := &Contract{
		Base: Base{
			Type:    SMART_CONTRACT,
			ChainId: chainId,
			Owner:   transaction.Owner,
			//Signers:  GetSigners(transaction.Owner),
			Signers:  transaction.Signers,
			Sequence: next.Sequence,
		},
		Data:     inputs,
		Function: COMPARE,
		Fee:      fee,
		Gas:      gas,
	}
	return SignAndPack(Transaction(compare))
}
