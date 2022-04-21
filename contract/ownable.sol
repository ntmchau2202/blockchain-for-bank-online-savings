// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./Context.sol";

abstract contract Ownable is Context {
    address private _owner;
    
    event OwnershipTransferred (
        address indexed _oldOwner,
        address indexed _newOwner
    );

     constructor() {
        address msgSender = _msgSender();
        _owner = msgSender;
        emit OwnershipTransferred(address(0), _owner);
    }

    function owner() external view returns (address) {
        return _owner;
    }

    modifier onlyOwner() {
        require(_owner == _msgSender(), "Error: caller is not owner");
        _;
    }

    function transferOwnership(address newOwner) public virtual onlyOwner {
        require(
            newOwner != address(0),
            "Error: new owner must not be zero address"
        );
        _owner = newOwner;
        emit OwnershipTransferred(_owner, newOwner);
    }

}

