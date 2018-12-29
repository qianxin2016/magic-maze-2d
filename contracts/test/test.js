var Test = artifacts.require("./MagicMaze.sol");
contract('MagicMaze', function(accounts) {
	it("call method create", function() {
	    Test.deployed().then(function(instance) {
		  return instance.call('create', 1, 'hello', 'ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad');
		}).then(function(result) {
          tronWeb.trx.getBalance(accounts[0]).then(result => { console.log(tronWeb.fromSun(result)); })
	    });
	});
	it("call method challenge", function() {
	    Test.deployed().then(function(instance) {
		  return instance.call('challenge', 1);
		}).then(function(result) {
          console.log(result);
	    });
	});
	it("call method takeBonus", function() {
	    Test.deployed().then(function(instance) {
		  return instance.call('takeBonus', 1, 'abc');
		}).then(function(result) {
          console.log(result);
	    });
	});
});



