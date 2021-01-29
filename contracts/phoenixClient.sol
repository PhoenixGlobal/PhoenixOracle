pragma solidity ^0.5.0;

import "./zeppelin/Ownable.sol";

contract phoenixClient is Ownable{
    bytes32 public value;
    uint public value1;
    uint public nonce;

    event Request(
        uint indexed nonce,
        address indexed to,
        bytes4 indexed fid
    );

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

    function requestData(address _address, bytes4 _fid) public{
        emit Request(nonce, _address, _fid);
        nonce += 1;
    }
}