package main

import (
  "time"
  "fmt"
  "./runner"
  "os"
)

func run(x chan string, y chan string) {
  runner := runner.CreateRunner()
  transaction, returnValue := runner.Call("0x0", `setWord('hello,world from Oneledger')`)
  x <- transaction
  y <- returnValue
}

func monitor(status_ch chan string) {
  i := 0
  threshold := 60
  for {
    time.Sleep(time.Second)
    i = i + 1
    if i >= threshold {
      fmt.Println("something is wrong")
      status_ch <- "crash the runtime"
      return
    }
  }
}

func main() {
  fmt.Println("starting OVM")
  defer func() {
    if r := recover(); r != nil {
      panic(r);
    }
    os.Remove("./ovm.pid")
  }()

  if _, err := os.Stat("./ovm.pid"); !os.IsNotExist(err) {
    panic("ovm.pid file exists, there is another ovm running or exit abnormally")
  }

  os.Create("./ovm.pid")

  transaction_ch := make(chan string)
  returnValue_ch := make(chan string)
  status_ch := make(chan string)
  go monitor(status_ch)
  go run(transaction_ch, returnValue_ch)
  ready := 0
  for {
    select {
    case transaction:= <- transaction_ch:
      fmt.Println(transaction)
      ready ++
    case returnValue := <- returnValue_ch:
      fmt.Println(returnValue)
      ready ++
    case status := <- status_ch:
      fmt.Println("retuning: ", status)
      panic("exit with code -1")
    }
    if ready == 2 {
      fmt.Println("ending OVM")
      return
    }
  }


}
