// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
import "./Ownable.sol";

contract BankFactory is Ownable {
    mapping (address => address) private listBankAddress;
    mapping (string => address) private listBank;

    constructor() {
        transferOwnership(msg.sender);
    }

    event NewBank(
        address indexed _bankOwner,
        address indexed _bankContractI,
        string indexed _bankIDI,
        string _bankName,
        string _bankID,
        address _bankContract,
        uint time
    );

    function registerNewBank (
        string memory _bankName,
        address _bankOwner,
        string memory _bankID
    ) external onlyOwner returns (address) {
        // require(listBankAddress[msg.sender].getFactory() == address(0x0), "Error: bank has already registered");
        require(listBank[_bankID] == address(0x0), "Error: there exists a bank with this bankID");
        Bank newBank = new Bank(
            _bankOwner,
            _bankName,
            _bankID
        );
        listBankAddress[_bankOwner] = address(newBank);
        listBank[_bankID] = _bankOwner;
        emit NewBank(_bankOwner, listBankAddress[_bankOwner], _bankID, _bankName, _bankID, listBankAddress[_bankOwner], block.timestamp);
        return listBankAddress[_bankOwner];
    }

    function getBankContractByID(
        string memory _bankID
    ) external view returns (address) {
        return listBankAddress[listBank[_bankID]];
    }

}

contract Bank is Ownable {
    struct User {
        address userAddress;
        string bankID;
    }

    // enum TransactionsType {
    //     OPEN_SAVINGS_ACCOUNT,
    //     SETTLE_SAVINGS_ACCOUNT
    // }

    mapping (address => User) private listUser;
    string private bankName;
    string private bankID;
    address private factory;


    constructor(
        address _bankOwner,
        string memory _bankName,
        string memory _bankID
    ) {
        transferOwnership(_bankOwner);
        bankName = _bankName;
        bankID = _bankID;
        factory = msg.sender;
    }

    function getBankName() external view returns (string memory) {
        return bankName;
    }

    function getBankID() external view returns (string memory) {
        return bankID;
    }

    event NewUser(
        address _userAddress,
        uint time
    );

    function getFactory() external view returns (address) {
        return factory;
    }

    function addUser (
        address _userAddress
    ) external onlyOwner {
        User memory newUser = User(_userAddress, bankID);
        listUser[_userAddress] = newUser;
        emit NewUser(_userAddress, block.timestamp);
    }

    function isMember(
        address _userAddress
    ) public view returns (bool) {
        User memory usr = listUser[_userAddress];
        if (usr.userAddress == address(0)) {
            return false;
        } else {
            return true;
        }
    }

    function verifyUser(
        bytes32 hash, 
        bytes[2] memory signature
    ) public view returns (bool, address) {
        bytes32 r;
        bytes32 s;
        uint8 v;
        address customerAddress = address(0x0);
        for (uint256 i = 0; i < signature.length; i++) {
            bytes memory currentSignature = signature[i];
            assembly {
                r := mload(add(currentSignature, 0x20))
                s := mload(add(currentSignature, 0x40))
                v := byte(0, mload(add(currentSignature, 0x60)))
            }

            // Version of signature should be 27 or 28, but 0 and 1 are also possible versions
            if (v < 27) {
                v += 27;
            }

            // If the version is correct return the signer address
            if (v != 27 && v != 28) {
                return (false, customerAddress);
            } else {
                address recovered = ecrecover(hash, v, r, s);
                User memory usr = listUser[recovered];
                if (!isMember(usr.userAddress) && recovered != this.owner()) {
                    return (false, customerAddress);
                }

                if (isMember(usr.userAddress)) {
                    customerAddress = usr.userAddress;
                }
            }
        }
        return (true, customerAddress);
    }

    event OpenSavingsAccountTransaction(
        address indexed _customerAddressI,
        address _customerAddress,
        bytes _firstSignature,
        bytes _secondSignature
    );

    event SettleSavingsAccountTransaction(
        address indexed _customerAddressI,
        address _customerAddress,
        bytes _firstSignature,
        bytes _secondSignature
    );

    function BroadcastOpenAccountTransaction(
        address _customerAddress,
        bytes32 _hash,
        bytes[2] memory _signatures
    ) public {
        bool valid = false;
        address customerAddress = address(0x0);
        (valid, customerAddress) = verifyUser(_hash, _signatures);
        require(valid, "Error: invalid signatures");
        require(customerAddress == _customerAddress, "Error: customer mismatch");
        emit OpenSavingsAccountTransaction(
            customerAddress,
            customerAddress,
            _signatures[0],
            _signatures[1]
        );
    }

    function BroadcastSettleAccountTransaction(
        address _customerAddress,
        bytes32 _hash,
        bytes[2] memory _signatures
    ) public {
        bool valid = false;
        address customerAddress = address(0x0);
        (valid, customerAddress) = verifyUser(_hash, _signatures);
        require(valid, "Error: invalid signatures");
        require(customerAddress == _customerAddress, "Error: customer mismatch");
        emit OpenSavingsAccountTransaction(
            customerAddress,
            customerAddress,
            _signatures[0],
            _signatures[1]
        );
    }
}