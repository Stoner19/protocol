var context = require('./context');
var Contract = require('./simple_contract');

var contract = new Contract(context);

contract.setWord('myaddress','hello,world')
var output = contract.getWord('myaddress')
console.log(output);
