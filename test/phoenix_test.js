const phoenixClient = artifacts.require("phoenixClient");
const web3 = require("web3")
contract('phoenixClient', () => {
    let oracle;
    beforeEach(async () => {
        oracle = await phoenixClient.deployed();
    });
    it("returns the value it was initialized with", async () => {
        let value1 = await oracle.value1.call();
        //let value = await oracle.value.call();
        assert.equal(1, value1);
        var value = web3.utils.toHex("hellp apex");
        console.log(value); // "1000000000000000000"
        //assert.equal(web3.utils.toHex("hellp apex"), web3.utils)
    });

    it("sets the value", async () => {
        await oracle.setValue1(2);
        let value1 = await oracle.value1.call();
        assert.equal(2, value1);
    });
});