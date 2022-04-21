// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./Ownable.sol";
import "./bank.sol";

// Model: one base contract
// banks want to connect to -> register to create contract for that own bank
// the bank request to contract -> contract call other contract -> return data

contract BankFactory is Ownable {
    mapping(string => address) private mapBankAddress;
    mapping(address => Bank) private mapBank;
    Bank[] listBank;
    Bank private newBank;

    event Open(
        string indexed _savingsAccountNumberI,
        string indexed _customerIDI,
        string indexed _bankNameI,
        string _savingsAccountNumber,
        string _customerID,
        string _productType,
        string _period,
        string _interestRate,
        string _savingsAmount,
        string _estimatedInterestAmount,
        string _settleInstruction,
        string _openTime,
        string _currency,
        string _bankName
    );

    event Settle(
        string indexed _savingsAccountNumberI,
        string indexed _customerIDI,
        string indexed _bankNameI,
        string _savingsAccountNumber,
        string _customerID,
        string _actualInterestRate,
        string _settleTime,
        string _bankName
    );

    event NewBank(
        string indexed _bankNameI,
        string _bankName,
        address _bankAddr
    );

    constructor () {
        transferOwnership(msg.sender);
    }

    function registerBank(
        string memory _bankName
    ) external onlyOwner returns (address) {
        require(mapBankAddress[_bankName] == address(0), "Bank has already joined the network");
        newBank = new Bank(_bankName);
        address addr = address(newBank);
        mapBankAddress[_bankName] = addr;
        mapBank[addr] = newBank;
        listBank.push( newBank);
        emit NewBank(_bankName, _bankName, addr);
        return addr;
    }

    function getBankAddress(
        string memory _bankName
    ) external view returns (address) {
        return mapBankAddress[_bankName];
    }

    function numOfBanks() external view returns (uint) {
        return listBank.length;
    }

    function getAllBanks() external view returns (Bank[] memory) {
        return listBank;
    }

    function invokeOpenAccount(
        string memory _bankName,
        string memory _savingsAccountNumber,
        string memory _customerID,
        string memory _productType,
        string memory _period,
        string memory _interestRate,
        string memory _savingsAmount,
        string memory _estimatedInterestAmount,
        string memory _settleInstruction,
        string memory _openTime,
        string memory _currency
    ) external {
        Bank targetBank = findBank(_bankName);
        targetBank.OpenSavingsAccount(
            _savingsAccountNumber,
            _customerID,
            _productType,
            _period,
            _interestRate,
            _savingsAmount,
            _estimatedInterestAmount,
            _settleInstruction,
            _openTime,
            _currency
        );

        emit Open(_savingsAccountNumber,
                    _customerID, 
                    _bankName,
                    _savingsAccountNumber,
                    _customerID,
                    _productType, 
                    _period, 
                    _interestRate, 
                    _savingsAmount,
                    _estimatedInterestAmount, 
                    _settleInstruction, 
                    _openTime, 
                    _currency,
                    _bankName);
    }

    function invokeSettleAccount(
        string memory _bankName,
        string memory _savingsAccountNumber,
        string memory _customerID,
        string memory _actualInterestRate,
        string memory _settleTime
    ) external {
        Bank targetBank = findBank(_bankName);
        targetBank.SettleSavingsAccount(_savingsAccountNumber, _customerID, _actualInterestRate, _settleTime);
        emit Settle(_savingsAccountNumber, 
                    _customerID,
                    _bankName,
                    _savingsAccountNumber, 
                    _customerID, 
                    _actualInterestRate,
                    _settleTime,
                    _bankName);
    }

    function findBank(
        string memory _bankName
    ) internal view returns (Bank) {
        address bankAddress = mapBankAddress[_bankName];
        require(bankAddress !=address(0) , "Bank does not join in chain");
        Bank targetBank = mapBank[bankAddress];
        return targetBank;
    }
}


