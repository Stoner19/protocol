package main

import (
  "fmt"
  "github.com/robertkrimen/otto"
)

func main() {
	fmt.Println("starting OVM")
  vm := otto.New()
  vm.Set("version", "OVM 0.1 simple test");
  vm.Run(`
    abc = 2+2;
    console.log("the value of abc is " + abc);
    console.log("running on the vm " + version);
    `)
}
