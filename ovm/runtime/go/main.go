package main

import (
  "time"
  "fmt"
  "sync"
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
    i = i + 1
    if i == 2 {
      panic("crash the runtime")
    }
  }
}


func main() {
  fmt.Println("starting OVM")
  transaction_ch := make(chan string)
  returnValue_ch := make(chan string)
  var done sync.WaitGroup
  go monitor()
  go run(transaction_ch, returnValue_ch)
  transaction, returnValue := <-transaction_ch, <-returnValue_ch
  fmt.Println(transaction)
  fmt.Println(returnValue)
  fmt.Println("ending OVM")
  done.Wait()
}
