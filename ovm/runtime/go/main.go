package main

import (
  "fmt"
  "github.com/robertkrimen/otto"
)

// this is going to get the code from block chain by the address string
func getCodeFromBlockChain(address string) string{
  return `
  var SimpleContract = function (context) {
    this.context = context;
  }

  SimpleContract.prototype.setWord = function(word) {
    this.context.set('word', word);
  }

  SimpleContract.prototype.getWord = function () {
    return this.context.get('word');
  }
  runner.Contract = SimpleContract
  `
}

func main() {
	fmt.Println("starting OVM")
  vm := otto.New()
  vm.Set("version", "OVM 0.1 TEST")
  //initialize context
  vm.Run(`
    var IndexSet = function () { ///this is a simple porting for the set which lacks from ES5
      this._innerArray = [];
    }

    IndexSet.prototype.add = function(val) {
      for (var i = 0; i < this._innerArray.length; i ++) {
        if (this._innerArray[i] === val) {
          return;
        }
      }
      this._innerArray.push(val);
    }

    IndexSet.prototype.remove = function(val) {
      var replacedArray = [];
      for (var i = 0; i < this._innerArray.length; i ++) {
        if (this._innerArray[i] != val) {
          replacedArray.push(this._innerArray[i]);
        }
      }
      this._innerArray = replacedArray;
    }

    IndexSet.prototype.has = function (val) {
      for (var i = 0; i < this._innerArray.length; i ++) {
        if (this._innerArray[i] === val) {
          return true;
        }
      }
      return false;
    }

    IndexSet.prototype.list = function () {
      return this._innerArray;
    }

    var Context = function () {
      this.storage = {};
      this.updateIndexSet = new IndexSet();
    }

    Context.prototype.get = function (key) {
      return this.storage[key];
    }

    Context.prototype.set = function (key, val) {
      this.updateIndexSet.add(key);
      this.storage[key] = val;
    }

    Context.prototype.getStorage = function () {
      return this.storage;
    }

    Context.prototype.getUpdateIndexList = function () {
      return this.updateIndexSet.list();
    }
    var context = new Context();
    var runner = {};
    `)
  vm.Run(`(function(runner){` + getCodeFromBlockChain("0x0000000000") + `})(runner)`)
  vm.Run(`
    var contract = new runner.Contract(context);
    contract.setWord('myaddress','hello,world');
    var list = context.getUpdateIndexList();
    var storage = context.getStorage();
    var transaction = {};
    for (var i = 0; i< list.length; i ++) {
      var key = list[i];
      transaction[key] = storage[key];
    }
    `)
  vm.Run(`var result = JSON.stringify(transaction)`);
  if value, err := vm.Get("result"); err == nil {
    fmt.Println(value)
  }
  fmt.Println("ending OVM")
}
