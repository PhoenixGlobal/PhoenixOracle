const testApi = require('./support/helpers.js')

contract('phoenixClient', () => {
    let phoenix = artifacts.require("../contracts/phoenixClient.sol");
    let GetterSetter = artifacts.require("../contracts/GetterSetter.sol");
    let oc;
    let fID = "0x12345678";
    let to = "0x80e29acb842498fe6591f020bd82766dce619d43";
    beforeEach(async () => {
        oc = await phoenix.new({from : oracle});
    });

    it("has a limited public interface", () => {
        testApi.checkPublicABI(phoenix, [
            "transferOwnership",
            "requestData",
            "fulfillData",
        ]);
    });

    describe("#transferOwnership", () => {
        context("when called by the owner", () => {
            beforeEach( async () => {
                await oc.transferOwnership(stranger, {from: oracle});
            });

            it("can change the owner", async () => {
                let owner = await oc.owner.call();
                assert.isTrue(web3.utils.isAddress(owner));
                assert.equal(stranger, owner);
            });
        });

        context("when called by a non-owner", () => {
            it("cannot change the owner", async () => {
                await testApi.assertActionThrows(async () => {
                    await oc.transferOwnership(stranger, {from: stranger});
                });
            });
        });
    });
});