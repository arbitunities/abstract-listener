package main

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/arbitunities/abstract-listener/internal/argparser"
	"github.com/arbitunities/abstract-listener/internal/logger"
	"github.com/arbitunities/abstract-listener/pkg/config"
	"github.com/arbitunities/abstract-listener/pkg/listener"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	topics := [][]common.Hash{}
	for _, event := range []string{
		"UserOperationEvent(bytes32,address,address,uint256,bool,uint256,uint256)",
		"UserOperationRevertReason(bytes32,address,uint256,bytes)",
		"AccountDeployed(bytes32,address,address,address)",
		"SignatureAggregatorChanged(address)",
	} {
		topics = append(topics, []common.Hash{
			crypto.Keccak256Hash([]byte(event)),
		})
	}

	// Go Ethererum light client connection
	client, err := ethclient.Dial(os.Getenv("RPC"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed create client %s\n", err)
		os.Exit(1)
	}

	if err := run(client, ethereum.FilterQuery{
		Addresses: []common.Address{config.EntryPointAddress},
		Topics:    topics,
		FromBlock: big.NewInt(16912244),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start listener %s\n", err)
		os.Exit(1)
	}
}

func run(client *ethclient.Client, filterQuery ethereum.FilterQuery) error {
	// parse flags
	debug, from, to := argparser.ParseFlags()

	// setup logger
	logger := logger.NewLogger(debug)

	// setup listener
	lol := listener.NewListener(client)
	txs, err := lol.GetAddressTransactions(
		context.Background(),
		&config.EntryPointAddress,
		big.NewInt(*from),
		big.NewInt(*to),
	)

	if err != nil {
		return err
	}

	for _, tx := range txs {
		logger.Log().
			Str("hash", tx.Hash).
			Int64("chain", tx.ChainId.Int64()).
			Str("from", tx.From).
			Str("to", tx.To).
			Uint64("nonce", tx.Nonce).
			Uint64("gas", tx.Gas).
			Uint64("gasprice", tx.GasPrice).
			// Str("data", tx.Data).
			Msg("")
	}
	return nil
}
