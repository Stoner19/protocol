/*
	Copyright 2017-2018 OneLedger

	An incoming transaction, send, swap, ready, verification, etc.
*/
package action

import (
	"github.com/Oneledger/protocol/node/err"
	"github.com/Oneledger/protocol/node/id"
	"github.com/tendermint/go-crypto"
)

type Message = []byte // Contents of a transaction
type PublicKey = crypto.PubKey

// ENUM for type
type Type byte
type Role byte

const (
	INVALID       Type = iota
	REGISTER           // Register a new identity with the chain
	SEND               // Do a normal send transaction on local chain
	EXTERNAL_SEND      // Do send on external chain
	EXTERNAL_LOCK      // Lock some data on external chain
	SWAP               // Start a swap between chains
	VERIFY             // Verify that a lockbox is correct
	PUBLISH            // Exchange data on a chain
	READ               // Read a specific transaction on a chain
	PREPARE            // Do everything, except commit
	COMMIT             // Commit to doing the work
	FORGET             // Rollback and forget that this happened
 	CHECKFORERROR	   //todo: check something happened(with a delay transaction), otherwise it's an error
)

const (
	ALL         Role = iota
	INITIATOR        // Register a new identity with the chain
	PARTICIPANT      // Do a normal send transaction on local chain
	NONE
)

// Polymorphism and Serializable
type Transaction interface {
	Validate() err.Code
	ProcessCheck(interface{}) err.Code
	ShouldProcess(interface{}) bool
	ProcessDeliver(interface{}) err.Code
	Expand(interface{}) Commands
	Resolve(interface{}, Commands)
}

// Base Data for each type
type Base struct {
	Type    Type   `json:"type"`
	ChainId string `json:"chain_id"`

	Owner   id.AccountKey `json:"owner"`
	Signers []PublicKey   `json:"signers"`

	Sequence int64 `json:"sequence"`
	Delay    int64 `json:"delay"` // Pause the transaction in the mempool
}
