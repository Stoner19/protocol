var Context = function () {
  this.storage = {};
}

Context.prototype.get = function (key) {
  return this.storage[key];
}

Context.prototype.set = function (key, val) {
  this.storage[key] = val;
}

Context.prototype.getStorage = function () {
  return this.storage;
}

module.exports = new Context();
