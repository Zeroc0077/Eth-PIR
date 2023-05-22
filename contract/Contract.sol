// SPDX-License-Identifier: GPL-3.0

pragma solidity ^ 0.8.7;

struct User {
    uint id;
    string identity;
}

struct Fair {
    uint id;
    uint pay;
}

contract Test {
    address private owner;  // owner of the contract
    mapping(address => User) private users;  // record the user's identity
    mapping(address => Fair) private pays;  // record the user's payment
    mapping(uint => address payable) private server_record;  // record the server's address
    mapping(uint => bool) private server_record_test;  // record the server's address
    uint count = 0;  // count the number of users
    bool isSuccess = false;  // record the result of the service
    uint servercount = 0;  // count the number of servers
    uint sum = 0;  // record the sum of the payment

    constructor() {
        owner = msg.sender;
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only Owner can modify the status");
        _;
    }

    modifier onlyClient() {
        require(keccak256(bytes(users[msg.sender].identity)) == keccak256(bytes("Client")), "Only Client can confirm the service");
        _;
    }

    modifier onlySuccess() {
        require(isSuccess == true, "The service result isn't confirmed by client");
        _;
    }

    // start the whole service process
    function startProcess() public payable {
        require(msg.value > 0, "Need some eths to start process");
        uint id = count;
        count += 1;
        users[msg.sender] = User(id, "Client");
        pays[msg.sender] = Fair(id, msg.value);
        sum += msg.value;
    }

    // charge the server
    function chargeServer() public payable {
        require(msg.value > 0, "Server need to pay some eths to serve client");
        uint id = count;
        count += 1;
        users[msg.sender] = User(id, "Server");
        pays[msg.sender] = Fair(id, msg.value);
        sum += msg.value;
        servercount += 1;
        server_record[id] = payable(msg.sender);
        server_record_test[id] = true;
    }

    // Initialize the server
    function isInitialized(uint key) private view returns (bool) {
        return server_record_test[key];
    }

    // if the service is success, pay the server
    function payServer() private onlySuccess {
        uint payoff = sum / servercount;
        for(uint i=0; i < count; ++i){
            if(isInitialized(i)){
                server_record[i].transfer(payoff);
            }
        }
    }

    // client confirm the service result
    function clientConfirm(bool isConfirm) public onlyClient {
        if(isConfirm){
            isSuccess = true;
            payServer();
        }
        else{
            isSuccess = false;
            selfDestruct();
        }
    }

    // if the service is failed, selfdestruct
    function selfDestruct() public {
        selfdestruct(payable(msg.sender));
    }
}