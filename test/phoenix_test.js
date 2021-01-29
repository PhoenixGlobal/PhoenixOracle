const phoenixClient = artifacts.require("phoenixClient");
// var Web3 = require('web3');
// var web3 = new Web3(new Web3.providers.HttpProvider("http://localhost:9545"))
contract('phoenixClient', () => {
    let oracle;
    beforeEach(async () => {
        oracle = await phoenixClient.deployed();
    });
    it("returns the value it was initialized with", async () => {
        let value = await oracle.value1.call();
        assert.equal(1, value);
    });
});