pragma solidity ^0.4.23;

contract MagicMaze {
    struct Maze {
        uint256 id;
        string name;
        bytes32 mazeHash;
        uint256 bonus;
        address creator;
        address challenger;
        address winner;
    }

    address public owner;
    uint256 constant SERVICE_FEE = 1 trx;
    mapping(uint256 => Maze) public mazes;
    uint256 public serviceCharge;

    event mazeCreated(uint256 id);
    event winner(uint256 id, address player);
    event falseClaim(uint256 id, address player);

    modifier onlyOwner() {
        require(msg.sender == owner);
        _;
    }

    constructor() public {
        owner = msg.sender;
    }

    function create(uint256 id, string memory name, string memory mazeHash) public payable {
        require(msg.value >= SERVICE_FEE);
        require(mazes[id].bonus == 0);
        bytes32 hash = stringToBytes32(mazeHash);
        mazes[id] = Maze(id, name, hash, msg.value, msg.sender, address(0), address(0));
        emit mazeCreated(id);
    }
    
    function challenge(uint256 id) public payable {
        require(mazes[id].bonus != 0);
        require(mazes[id].challenger == address(0));
        require(mazes[id].bonus == msg.value);
        
        mazes[id].challenger = msg.sender;
        mazes[id].bonus += msg.value;
    }

    function takeBonus(uint256 id, string memory mazeInfo) public payable {
        require(mazes[id].bonus != 0);
        require(mazes[id].winner == address(0));
        require(mazes[id].creator == msg.sender || mazes[id].challenger == msg.sender);
        require(msg.value >= SERVICE_FEE);

        bytes32 mazeHash = sha256(abi.encodePacked(mazeInfo));
        if (mazeHash == mazes[id].mazeHash) {
            mazes[id].winner = msg.sender;
            msg.sender.transfer(mazes[id].bonus);
            emit winner(id, msg.sender);
        } else {
            emit falseClaim(id, msg.sender);
        }

        serviceCharge += msg.value;
    }

    function getPlayerInfo(uint256 id) public view returns (address, address, address) {
        return (mazes[id].creator, mazes[id].challenger, mazes[id].winner);
    }

    function withdrawServiceCharge() public onlyOwner {
        owner.transfer(serviceCharge);
    }

    function uintToAscii(uint number) private pure returns (byte) {
        if (number < 10) {
            return byte(uint8(48 + number));
        } else if (number < 16) {
            return byte(uint8(87 + number));
        } else {
            revert();
        }
    }

    function asciiToUint(byte char) private pure returns (uint) {
        uint8 asciiNum = uint8(char);
        if (asciiNum > 47 && asciiNum < 58) {
            return asciiNum - 48;
        } else if (asciiNum > 96 && asciiNum < 103) {
            return asciiNum - 87;
        } else {
            revert();
        }
    }

    function bytes32ToString (bytes32 data) private pure returns (string memory) {
        bytes memory bytesString = new bytes(64);
        for (uint j=0; j < 32; j++) {
            uint char = uint(data) * 2 ** (8 * j);
            bytesString[j*2+0] = uintToAscii(char / 16);
            bytesString[j*2+1] = uintToAscii(char % 16);
        }
        return string(bytesString);
    }

    function stringToBytes32(string memory str) private pure returns (bytes32) {
        bytes memory bString = bytes(str);
        uint uintString;
        if (bString.length != 64) {
            revert();
        }
        for (uint i = 0; i < 64; i++) {
            uintString = uintString*16 + asciiToUint(bString[i]);
        }
        return bytes32(uintString);
    }
}
