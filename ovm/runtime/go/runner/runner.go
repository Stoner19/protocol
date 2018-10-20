package runner

import (
  "github.com/robertkrimen/otto"
)

type Runner struct {
  vm *otto.Otto
}

func (runner Runner) getContract(address string) bool {
  code :=  `
  var SimpleContract = function (context) {
    this.context = context;
  }

  SimpleContract.prototype.setWord = function(word) {
    this.context.set('word', word);
    return word;
  }

  SimpleContract.prototype.getWord = function () {
    return this.context.get('word');
  }
  module.Contract = SimpleContract;
  `
  _, error := runner.vm.Run(`var module = {};(function(module){` + code + `})(module)`)
  if error == nil {
    return true
  } else {
    return false
  }
}

func (runner Runner) initialContext() {
  runner.vm.Run(`
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
    `)
}

func (runner Runner) exec(callString string) (string, string) {
  _, error := runner.vm.Run(`
    var a = 0;
    while(true){
      a += 1;
      if (a > 1000000){
        a = 0;
        console.log('dead_loop');
      }

    }
    var contract = new module.Contract(context);
    var retValue = contract.` + callString)
  if error != nil {
    panic(error)
  }
  runner.vm.Run(`
    var list = context.getUpdateIndexList();
    var storage = context.getStorage();
    var transaction = {};
    for (var i = 0; i< list.length; i ++) {
      var key = list[i];
      transaction[key] = storage[key];
    }
    `)
  runner.vm.Run(`
    transaction = JSON.stringify(transaction);
    retValue = JSON.stringify(retValue);
    `);
  output := ""
  returnValue := ""

  if value, err := runner.vm.Get("transaction"); err == nil {
    output, _ = value.ToString()
  }

  if value, err := runner.vm.Get("retValue"); err == nil {
    returnValue, _ = value.ToString()
  }
  return output, returnValue
}

func (runner Runner) Call(address string, callString string) (string, string){
  runner.initialContext()
  runner.getContract(address)
  transaction, returnValue := runner.exec(callString)
  return transaction, returnValue
}

func CreateRunner() Runner{
  vm := otto.New()
  vm.Set("version", "OVM 0.1 TEST")
  return Runner {vm}
}
