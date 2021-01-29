const Migrations = artifacts.require("Migrations");
const phoenixClient = artifacts.require("phoenixClient");

module.exports = function (deployer) {
  deployer.deploy(Migrations);
  deployer.deploy(phoenixClient);
};
