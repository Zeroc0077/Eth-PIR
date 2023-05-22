// Package Utils
/**
 * @author zeroc
 * @date 1:38 2023/5/21
 * @file Utils.go
 **/
package Utils

import (
	"Eth-PIR/contract"
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

/*
the ethclient is a client to the Ethereum client which is based on the JSON-RPC API.
It provides the following methods:
	- BlockByHash
	- BlockByNumber
	- BlockNumber
	- CallContract
	- ChainID
	- CodeAt
	- EstimateGas
	- FilterLogs
	- HeaderByHash
	- HeaderByNumber
	- NetworkID
	- PendingBalanceAt
	- PendingCallContract
	- PendingCodeAt
	- PendingNonceAt
	- PendingTransactionCount
	- SuggestGasPrice
	- TransactionByHash
	- TransactionCount
	- TransactionInBlock
	- TransactionReceipt
	- SendTransaction
	- SendRawTransaction
	- Subscribe
	- SubscribeNewHead
	- SubscribeNewPendingTransactions
	- SubscribeLogs
	- SubscribeSynced
	- SyncProgress
	- TransactionSender
	- BalanceAt
	- FilterLogs
*/

// Connect to the Ethereum client
func Connect(url string) *ethclient.Client {
	// Dial is used to create a client connection to the given URL.
	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	return client
}

// GetLatestBlockNumber Get the latest block number
func GetLatestBlockNumber(client *ethclient.Client) uint64 {
	// context is the execution context for the call, carrying deadline, cancellation and other values across API boundaries.
	// Background returns a non-nil, empty Context. It is never canceled, has no values, and has no deadline.
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatalf("Failed to retrieve latest block number: %v", err)
	}
	return blockNumber
}

// GetBlockByNumber Get block information by block number
func GetBlockByNumber(client *ethclient.Client, blockNumber uint64) *types.Block {
	// BlockByNumber returns the given full block. If number is nil, the latest known block is returned.
	var block *types.Block
	block, err := client.BlockByNumber(context.Background(), big.NewInt(int64(blockNumber)))
	if err != nil {
		log.Fatalf("Failed to retrieve block information: %v", err)
	}
	return block
}

// GetTransactionByHash Get transaction information by transaction hash
func GetTransactionByHash(client *ethclient.Client, txHash string) *types.Transaction {
	// TransactionByHash returns the transaction with the given hash.
	var tx *types.Transaction
	tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	if err != nil {
		log.Fatalf("Failed to retrieve transaction information: %v", err)
	}
	return tx
}

// GetTransactionReceipt Get transaction receipt by transaction hash
func GetTransactionReceipt(client *ethclient.Client, txHash string) *types.Receipt {
	// TransactionReceipt returns the receipt of a transaction by transaction hash.
	var txReceipt *types.Receipt
	txReceipt, err := client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		log.Fatalf("Failed to retrieve transaction receipt: %v", err)
	}
	return txReceipt
}

// Close the connection
func Close(client *ethclient.Client) {
	client.Close()
}

// GetBalance Get balance by address
func GetBalance(client *ethclient.Client, addresshex string) (*big.Float, error) {
	// BalanceAt returns the wei balance of the given account at the given block.
	address := common.HexToAddress(addresshex)
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, err
	}
	// transfer balance to ether
	ethBalance := WeitoEth(big.NewFloat(0).SetInt(balance))

	return ethBalance, nil
}

// GetAddress Get address by private key
func GetAddress(privatekey string) string {
	// HexToECDSA parses a secp256k1 private key.
	key, err := crypto.HexToECDSA(privatekey)
	if err != nil {
		log.Fatalf("Failed to get address: %v", err)
	}
	// Public returns the ECDSA public key corresponding to this private key.
	publickey := key.Public()
	publickeyECDSA, ok := publickey.(*ecdsa.PublicKey)
	if !ok {
		panic("cannot assert type: publickey is not of type *ecdsa.PublicKey")
	}
	// PubkeyToAddress returns the address corresponding to the given ECDSA public key.
	address := crypto.PubkeyToAddress(*publickeyECDSA).Hex()
	return address
}

// WeitoEth Transfer wei to eth
func WeitoEth(wei *big.Float) *big.Float {
	// 1 ether = 1e18 wei
	ethBalance := new(big.Float).Quo(wei, big.NewFloat(1e18))
	return ethBalance
}

// StartProcess Start the process
func StartProcess(privatekey string, chainID *big.Int) {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	cont, err := contract.NewContract(common.HexToAddress("0xc063eB8efC5EE6aada0b34B68e09e469921052b6"), client)
	key, err := crypto.HexToECDSA(privatekey)
	if err != nil {
		log.Fatalf("Failed to get key: %v", err)
	}
	// NewKeyedTransactor creates a new transactor from a secp256k1 private key and chainID.
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Fatalf("Failed to get auth: %v", err)
	}
	auth.Value = big.NewInt(1e18)
	// the method of a contract needs the auth which provides the account to send the transaction from.
	_, err = cont.StartProcess(auth)
	if err != nil {
		log.Fatalf("Failed to start process: %v", err)
	}
}

// ChargeServer Charge the server
func ChargeServer(privatekey string, chainID *big.Int) {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	cont, err := contract.NewContract(common.HexToAddress("0xc063eB8efC5EE6aada0b34B68e09e469921052b6"), client)
	key, err := crypto.HexToECDSA(privatekey)
	if err != nil {
		log.Fatalf("Failed to get key: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Fatalf("Failed to get auth: %v", err)
	}
	auth.Value = big.NewInt(1e18)
	_, err = cont.ChargeServer(auth)
	if err != nil {
		log.Fatalf("Failed to charge server: %v", err)
	}
}

// ClientConfirm Client confirm
func ClientConfirm(privatekey string, chainID *big.Int, isConfirm bool) {
	client, err := ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	cont, err := contract.NewContract(common.HexToAddress("0xc063eB8efC5EE6aada0b34B68e09e469921052b6"), client)
	key, err := crypto.HexToECDSA(privatekey)
	if err != nil {
		log.Fatalf("Failed to get key: %v", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Fatalf("Failed to get auth: %v", err)
	}
	_, err = cont.ClientConfirm(auth, isConfirm)
	if err != nil {
		log.Fatalf("Failed to client confirm: %v", err)
	}
}
