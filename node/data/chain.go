/*
	Copyright 2017 - 2018 OneLedger

	Define chains as they are seen in OneLedher
*/
package data

import "github.com/Oneledger/protocol/node/serial"

type ChainType int

// TODO: These should be in a domain database
const (
	UNKNOWN ChainType = iota
	ONELEDGER
	BITCOIN
	ETHEREUM
)

func (ctype ChainType) String() string {
	switch ctype {
	case UNKNOWN:
		return "UNKNOWN"
	case ONELEDGER:
		return "OneLedger"
	case BITCOIN:
		return "BitCoin"
	case ETHEREUM:
		return "Ethereum"
	default:
		return "INVALID"
	}
}

type Chain struct {
	ChainType   ChainType
	Description string
	Features    []string
}

// A specific node in a chain
type ChainNode struct {
	ChainType ChainType
	Location  string

	// TODO: Causing cycle...
	//Owner     id.Identity
}

func init() {
	serial.Register(ChainType(0))
	serial.Register(Chain{})
	serial.Register(ChainNode{})
}
