/*
	Copyright 2017 - 2018 OneLedger

	Handle setting any options for the node.
*/
package app

import (
	"github.com/Oneledger/protocol/node/comm"
	"github.com/Oneledger/protocol/node/id"
	"github.com/Oneledger/protocol/node/log"
)

// Arguments for registration
type RegisterArguments struct {
	Identity   string
	Chain      string
	PublicKey  string
	PrivateKey string
}

func SetOption(app *Application, key string, value string) bool {
	log.Debug("Setting Application Options", "key", key, "value", value)

	switch key {

	case "Register":
		var arguments RegisterArguments
		result, err := comm.Deserialize([]byte(value), &arguments)
		if err != nil {
			log.Error("Can't set options", "err", err)
			return false
		}
		args := result.(*RegisterArguments)
		publicKey, privateKey := id.GenerateKeys([]byte(args.Identity)) // TODO: Switch with passphrase
		RegisterLocally(app, args.Identity, "OneLedger", id.ParseAccountType(args.Chain),
			publicKey, privateKey)

	default:
		log.Warn("Unknown Option", "key", key)
		return false
	}
	return true

}
