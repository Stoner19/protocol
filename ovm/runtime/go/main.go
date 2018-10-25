package main

import (
  "fmt"
  "os"
  "./runner"
  "./monitor"
)

func run(x chan string, y chan string) {
  runner := runner.CreateRunner()
  transaction, returnValue := runner.Call("0x0", `setWord('hello,world from Oneledger')`)
  x <- transaction
  y <- returnValue
}

func main() {
  fmt.Println("starting OVM")
  transaction_ch := make(chan string)
  returnValue_ch := make(chan string)
  status_ch := make(chan monitor.Status)
  monitor := monitor.CreateMonitor(10, monitor.DEFAULT_MODE, "./ovm.pid")

  status, err := monitor.CheckUnique()

  defer func() {
    if r := recover(); r != nil {
      fmt.Println(status)
      panic(r);
    }
    os.Remove(monitor.GetPidFilePath())
  }()

  if err == true  {
    panic(status)
  } else {
    fmt.Println(status)
  }

  os.Create(monitor.GetPidFilePath())


  go monitor.CheckStatus(status_ch)
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
      fmt.Println("retuning: ", status.Details, "with code", status.Code)
      panic("exit with code -1")
    }
    if ready == 2 {
      fmt.Println("ending OVM")
      return
    }
  }


}
