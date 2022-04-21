
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Bank {
    string private bankName;
    string private bankID;
    constructor (
        string memory _bankName
    ) {
        bankName = _bankName;
    }

    function getBankAddress() external view returns (address) {
        return address(this);
    }

    function getBankName() external view returns (string memory)  {
        return bankName;
    }   

    event openSavingsAccount(
        string indexed _savingsAccountNumberI,
        string indexed _customerIDI,
        string _savingsAccountNumber,
        string _customerID,
        string _productType,
        string _period,
        string _interestRate,
        string _savingsAmount,
        string _estimatedInterestAmount,
        string _settleInstruction,
        string _openTime,
        string _currency
    );

    event settleSavingsAccount(
        string indexed _savingsAccountNumberI,
        string indexed _customerIDI,
        string _savingsAccountNumber,
        string _customerID,
        string _actualInterestRate,
        string _settleTime
    );

    function OpenSavingsAccount(
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
        emit openSavingsAccount(
            _savingsAccountNumber,
            _customerID,
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
    }

    function SettleSavingsAccount (
        string memory _savingsAccountNumber,
        string memory _customerID,
        string memory _actualInterestRate,
        string memory _settleTime
    ) external {
        emit settleSavingsAccount(
            _savingsAccountNumber,
            _customerID,
            _savingsAccountNumber,
            _customerID,
            _actualInterestRate,
            _settleTime
        );
    }
}
