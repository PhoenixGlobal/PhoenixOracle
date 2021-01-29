BigNumber = require('bignumber.js');
moment = require('moment');
// const web3 = require("web3")
var Web3 = require('web3');
var web3 = new Web3(new Web3.providers.HttpProvider("http://localhost:9545"))
const eth = web3.eth;

// before(async function () {
//     var accounts = await eth.accounts;
//     // var Accounts = accounts.slice(1);
// });
//
// var Eth = function sendEth(method, params) {
//     params = params || [];
//
//     return new Promise((resolve, reject) => {
//         web3.currentProvider.sendAsync({
//             jsonrpc: "2.0",
//             method: method,
//             params: params || [],
//             id: new Date().getTime()
//         }, function sendEthResponse(error, response) {
//             if (error) {
//                 reject(error);
//             } else {
//                 resolve(response.result);
//             };
//         }, () => {}, () => {});
//     });
// };
//
//   const emptyAddress = '0x0000000000000000000000000000000000000000';
//   async function sealBlock() {
//     return Eth('evm_mine');
//   };
//
//   async function sendTransaction(params) {
//     return await eth.sendTransaction(params);
//   }
//
//   async function getBalance(account) {
//     return bigNum(await eth.getBalance(account));
//   }
//
//   function bigNum(number) {
//     return new BigNumber(number);
//   }
//
//   function toWei(number) {
//     return bigNum(web3.toWei(number));
//   }
//
//   function tokens(number) {
//     return bigNum(number * 10**18);
//   }
//
//   function intToHex(number) {
//     return '0x' + bigNum(number).toString(16);
//   }
//
//   function hexToInt(string) {
//     return web3.toBigNumber(string);
//   }
//
//   function hexToAddress(string) {
//     return '0x' + string.slice(string.length - 40);
//   }
//
//   function unixTime(time) {
//     return moment(time).unix();
//   }

  exports.seconds = function seconds(number) {
    return number;
  };

  // function minutes(number) {
  //   return number * 60;
  // };
  //
  // function hours(number) {
  //   return number * minutes(60);
  // };
  //
  // function days(number) {
  //   return number * hours(24);
  // };
  //
  // function keccak256(string) {
  //   return web3.sha3(string);
  // }
  //
  // function logTopic(string) {
  //   let hash = keccak256(string);
  //   return '0x' + hash.slice(26);
  // }
  //
  // async function getLatestBlock() {
  //   return await eth.getBlock('latest', false);
  // };
  //
  // async function getLatestTimestamp () {
  //   let latestBlock = await getLatestBlock()
  //   return web3.toDecimal(latestBlock.timestamp);
  // };
  //
  // async function fastForwardTo(target) {
  //   let now = await getLatestTimestamp();
  //   assert.isAbove(target, now, "Cannot fast forward to the past");
  //   let difference = target - now;
  //   await Eth("evm_increaseTime", [difference]);
  //   await sealBlock();
  // };
  //
  // function getEvents(contract) {
  //   return new Promise((resolve, reject) => {
  //     contract.allEvents().get((error, events) => {
  //       if (error) {
  //         reject(error);
  //       } else {
  //         console.log(events)
  //         resolve(events);
  //       };
  //     });
  //   });
  // };
  //
  // function eventsOfType(events, type) {
  //   let filteredEvents = [];
  //   for (event of events) {
  //     if (event.event === type) filteredEvents.push(event);
  //   }
  //   return filteredEvents;
  // };
  //
  // async function getEventsOfType(contract, type) {
  //   return eventsOfType(await getEvents(contract), type);
  // };
  //
  // async function getLatestEvent(contract) {
  //   let events = await getEvents(contract);
  //   return events[events.length - 1];
  // };
  //
  // function assertActionThrows(action) {
  //   return Promise.resolve().then(action)
  //     .catch(error => {
  //       assert(error, "Expected an error to be raised");
  //       assert(error.message, "Expected an error to be raised");
  //       return error.message;
  //     })
  //     .then(errorMessage => {
  //       assert(errorMessage, "Expected an error to be raised");
  //       assert.include(errorMessage, "invalid opcode", 'expected error message to include "invalid JUMP"');
  //       // see https://github.com/ethereumjs/testrpc/issues/39
  //       // for why the "invalid JUMP" is the throw related error when using TestRPC
  //     })
  // };
  //
  // function encodeUint256(int) {
  //   let zeros = "0000000000000000000000000000000000000000000000000000000000000000";
  //   let payload = int.toString(16);
  //   return (zeros + payload).slice(payload.length);
  // }
  //
  // function encodeAddress(address) {
  //   return '000000000000000000000000' + address.slice(2);
  // }
  //
  // function encodeBytes(bytes) {
  //   let zeros = "0000000000000000000000000000000000000000000000000000000000000000";
  //   let padded = bytes.padEnd(64, 0);
  //   let length = encodeUint256(bytes.length / 2);
  //   return length + padded;
  // }
  //
  // function checkPublicABI(contract, expectedPublic) {
  //   let actualPublic = [];
  //   for (method of contract.abi) {
  //     if (method.type == 'function') actualPublic.push(method.name);
  //   };
  //
  //   for (method of actualPublic) {
  //     let index = expectedPublic.indexOf(method);
  //     assert.isAtLeast(index, 0, (`#${method} is NOT expected to be public`))
  //   }
  //
  //   for (method of expectedPublic) {
  //     let index = actualPublic.indexOf(method);
  //     assert.isAtLeast(index, 0, (`#${method} is expected to be public`))
  //   }
  // };
  //
  // function functionID(signature) {
  //   return web3.sha3(signature).slice(2).slice(0, 8);
  // };
