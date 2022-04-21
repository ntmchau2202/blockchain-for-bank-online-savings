// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SavingsAccount {
    struct SavingsAccountDetails {
        string bankAccountNumber;
        address bank; // for the future, when we implement 1 bank => 1 blockchain account
        string savingsAccountNumber;
        string timeCreated;
        string interestRate; // how should we address the problem of float?
        string savingType;
        string savingPeriod; // months
        string savingAmount;
        string transactionUnit;
        string interestAmount;
        string totalReturnAmount;
        string timeSettled;
    }

    struct SavingsAccountTxn {
        string transactionType;
        string dateIssued;
        SavingsAccountDetails accountDetails;
    }

    string constant openSavingsAccCmd = "CREATE_SAVINGS_ACCOUNT";
    string constant settleSavingsAccCmd = "SETTLE_SAVINGS_ACCOUNT";

    // maps a bank account to a map
    // this map maps a saving account to transactions
    mapping(string => mapping(string => SavingsAccountTxn[])) savingsAccountTxns; 

    // event OpenSavingsAccount (
    //     string indexed _bankAccountNumber,
    //     string indexed _savingsAccountNumber,
    //     string indexed _timeIssued,
    //     string _timeCreated,
    //     string _savingAmount,
    //     string _savingPeriod,
    //     string _interestRate,
    //     string _savingType,
    //     string _transactionUnit
    // );

    event OpenSavingsAccount(
        string indexed _bankAccountNumber,
        string indexed _savingsAccountNumber,
        SavingsAccountTxn _openSavingsAccTxn);

    // event SettleSavingAccount (
    //     string indexed _bankAccountNumber,
    //     string indexed _savingsAccountNumber,
    //     string indexed _timeIssued,
    //     string _timeSettled,
    //     string _interestAmount,
    //     string _totalAmount,
    //     string _transactionUnit
    // );

    event SettleSavingsAccount(
            string indexed _bankAccountNumber,
            string indexed _savingsAccountNumber,
            SavingsAccountTxn _settleSavingsAccTxn
        );

    function registerSavingsAccount(
        string memory _bankAccountNumber,
        string memory _savingsAccountNumber,
        string memory _timeIssued,
        string memory _timeCreated,
        string memory _savingAmount,
        string memory _savingPeriod,
        string memory _interestRate,
        string memory _savingType,
        string memory _transactionUnit
    ) public {
        SavingsAccountDetails memory details = SavingsAccountDetails(
            _bankAccountNumber,
            msg.sender, 
            _savingsAccountNumber,
            _timeCreated,
            _interestRate, 
            _savingType,
            _savingPeriod, 
            _savingAmount,
            _transactionUnit,
            "0",
            "0",
            ""
        );

        SavingsAccountTxn memory transaction = SavingsAccountTxn(
            openSavingsAccCmd,
            _timeIssued,
            details
        );

        savingsAccountTxns[_bankAccountNumber][_savingsAccountNumber].push(transaction);
        emit OpenSavingsAccount(_bankAccountNumber, _savingsAccountNumber, transaction);
    }

    function settleSavingsAccount(
        string memory _bankAccountNumber,
        string memory _savingsAccountNumber,
        string memory  _timeIssued,
        string memory  _timeSettled,
        string memory  _interestAmount,
        string memory  _totalAmount
    ) public {
        require(keccak256(bytes(savingsAccountTxns[_bankAccountNumber][_savingsAccountNumber][0].transactionType)) != keccak256(bytes("")), "Savings account does not exist");
        SavingsAccountTxn memory transaction = savingsAccountTxns[_bankAccountNumber][_savingsAccountNumber][0];

        // create new transaction
        transaction.dateIssued = _timeIssued;
        transaction.transactionType = settleSavingsAccCmd;
        transaction.accountDetails.timeSettled = _timeSettled;
        transaction.accountDetails.interestAmount = _interestAmount;
        transaction.accountDetails.totalReturnAmount = _totalAmount;

        savingsAccountTxns[_bankAccountNumber][_savingsAccountNumber].push(transaction);
        emit SettleSavingsAccount(_bankAccountNumber, _savingsAccountNumber, transaction);
    }
}