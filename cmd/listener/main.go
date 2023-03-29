package main

import (
	"context"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/arbitunities/abstract-listener/pkg/config"
	"github.com/arbitunities/abstract-listener/pkg/listener"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
)

var (
	entryPointAddress    = common.HexToAddress("0x0576a174D229E3cFA37253523E645A78A0C91B57")
	simpleAccountFactory = common.HexToAddress("0x71D63edCdA95C61D6235552b5Bc74E32d8e2527B")
	addresses            = []common.Address{entryPointAddress, simpleAccountFactory}
)

func main() {
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

	config := config.NewConfig()
	// Go Ethererum light client connection
	client, err := ethclient.Dial(config.Rpc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed create client %s\n", err)
		os.Exit(1)
	}

	if err := run(client, ethereum.FilterQuery{
		Addresses: addresses,
		Topics:    topics,
		FromBlock: big.NewInt(16912244),
	}); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start listener %s\n", err)
		os.Exit(1)
	}
}

func run(client *ethclient.Client, filterQuery ethereum.FilterQuery) error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	logger := zerolog.New(os.Stdout)

	debug := flag.Bool("debug", false, "Debug log level ")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug || (strings.ToUpper(os.Getenv("LOG_LEVEL")) == "DEBUG") {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	lol := listener.NewListener(client)
	txs, err := lol.GetAddressTransactions(context.Background(), &entryPointAddress, int(16903038), int(16903038))

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
