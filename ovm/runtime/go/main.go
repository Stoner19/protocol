package main

import (
  "time"
  "fmt"
  "./runner"
)

func run(x chan string, y chan string) {

  runner := runner.CreateRunner()
  transaction, returnValue := runner.Call("0x0", `setWord('hello,world from Oneledger')`)
  x <- transaction
  y <- returnValue
}

func monitor() {
  i := 0
  for {
    time.Sleep(time.Second)
    fmt.Println(i)
    i = i + 1
  }
}


func main() {
  fmt.Println("starting OVM")
  transaction_ch := make(chan string)
  returnValue_ch := make(chan string)
  go monitor()
  go run(transaction_ch, returnValue_ch)
  transaction, returnValue := <-transaction_ch, <-returnValue_ch
  fmt.Println(transaction)
  fmt.Println(returnValue)
  fmt.Println("ending OVM")
  for {

  }
}
