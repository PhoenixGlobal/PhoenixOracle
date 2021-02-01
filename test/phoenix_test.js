const testApi = require('./support/helpers.js')

contract('phoenixClient', () => {
    let phoenix = artifacts.require("../contracts/phoenixClient.sol");

    let oracle;
    beforeEach(async () => {
        oracle = await phoenix.new();
    });
    describe("initialization", () => {
        it("returns the value it was initialized with", async () => {
            let value1 = await oracle.value1.call();
            //let value = await oracle.value.call();
            assert.equal(1, value1);
            var value = web3.utils.toHex("hellp apex");
            console.log(value); // "1000000000000000000"
            console.log(testApi.seconds(12));
            console.log(testApi.emptyAddress);
            //assert.equal(web3.utils.toHex("hellp apex"), web3.utils)
        });
    });

    describe("initialization", () => {
        it("sets the value", async () => {
            await oracle.setValue1(2);
            let value1 = await oracle.value1.call();
            assert.equal(2, value1);
        });
    });
});