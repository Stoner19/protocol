package main

import (
  "fmt"
  "github.com/robertkrimen/otto"
)

// this is going to get the code from block chain by the address string
func getCodeFromBlockChain(address string) string{
  return ``；
}

func main() {
	fmt.Println("starting OVM")
  vm := otto.New()
  vm.Set("version", "OVM 0.1 TEST")
  //initialize context
  vm.Run(``)
  vm.Run(getCodeFromBlockChain("0x0000000000"))
  fmt.Println("ending OVM")
}
