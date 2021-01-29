pragma solidity ^0.5.16;

contract phoenixClient {
    bytes32 public value;
    uint public value1;

    constructor () public{
        value = "hello apex";
        value1 = 1;
    }

    function setValue(bytes32 newValue) public {
        value = newValue;
    }

    function setValue1(uint newValue) public {
        value1 = newValue;
    }
}
