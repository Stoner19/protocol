var SimpleContract = function (context) {
  this.context = context;
}

SimpleContract.prototype.default__ = function() {
  this.context.set('word', 'hello, by default');
  return 'hello, by default';
}

SimpleContract.prototype.setWord = function(word) {
  this.context.set('word', word);
  return word;
}

SimpleContract.prototype.getWord = function () {
  return this.context.get('word');
}
module.Contract = SimpleContract;
