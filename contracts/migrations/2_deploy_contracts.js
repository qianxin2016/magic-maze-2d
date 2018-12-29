var Migrations = artifacts.require("./Migrations.sol");
var MagicMaze = artifacts.require("./MagicMaze.sol");
module.exports = function(deployer) {
  deployer.deploy(Migrations);
  deployer.deploy(MagicMaze);
};
