var SimpleContract = function (context) {
  this.context = context;
}


SimpleContract.prototype.default__ = function () {
  while(true){
    console.log('loop...');
  }
}
module.Contract = SimpleContract;
