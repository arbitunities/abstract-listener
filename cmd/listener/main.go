package main

import (
	"fmt"
	"os"

	"github.com/arbitunities/abstract-listener/pkg/listener"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start %s\n", err)
		os.Exit(1)
	}
}

func run() error {
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

	err := listener.NewListener().Run(topics)
	if err != nil {
		return err
	}
	return nil
}
