package main

import (
  "fmt"
  "./runner"
)

func main() {
	fmt.Println("starting OVM")
  runner := runner.CreateRunner()
  transaction, returnValue := runner.Call("0x0", `setWord('hello,world from Oneledger')`)
  fmt.Println(transaction)
  fmt.Println(returnValue)
  fmt.Println("ending OVM")
}
