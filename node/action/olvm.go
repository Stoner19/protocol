/*
	Copyright 2017-2018 OneLedger
*/

package action

import (
	"github.com/Oneledger/protocol/node/serial"
	"github.com/Oneledger/protocol/node/status"
)

type OLVMContext struct {
	Data interface{}
}

// All of the input necessary to perform a computation on a transaction
type OLVMRequest struct {
	From        string
	Address     string
	CallString  string
	Value       int
	SourceCode  []byte
	Transaction Transaction
	Context     OLVMContext

	// TODO: Data Handle (some way to call out for large data requests)
}

// All of the output received from the computation
type OLVMResult struct {
	Status  status.Code
	Out     string
	Ret     string // TODO: Should be a real name
	Elapsed string

	Transactions []Transaction
	Context      OLVMContext
}

func init() {
	serial.Register(OLVMRequest{})
	serial.Register(OLVMResult{})
	serial.Register(OLVMContext{})

	// TODO: Doesn't work in serial?
	//var prototype time.Time
	//serial.Register(prototype)
	//var prototype2 time.Duration
	//serial.Register(prototype2)
}

func NewOLVMRequest(script []byte, context OLVMContext) *OLVMRequest {
	request := &OLVMRequest{
		From:       "0x0",
		Address:    "embed://",
		CallString: "",
		Value:      0,
		SourceCode: script,
		Context:    context,
	}
	return request
}

func NewOLVMResult() *OLVMResult {
	result := &OLVMResult{
		Status: status.MISSING_DATA,
	}
	return result
}
