package listener

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/arbitunities/abstract-listener/pkg/chains"
	"github.com/arbitunities/abstract-listener/pkg/config"
	"github.com/common-nighthawk/go-figure"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Listener struct{}

func NewListener() *Listener {
	return &Listener{}
}

func (l *Listener) Run(topics [][]common.Hash) error {
	config := config.NewConfig()
	// Go Ethererum light client connection
	client, err := ethclient.Dial(config.Rpc)
	if err != nil {
		return err
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	chainId, err := client.ChainID(ctx)
	if err != nil {
		cancel()
		return err
	}

	fig := figure.NewFigure(chains.NameById(chainId), "starwars", true)
	fig.Print()

	headers := make(chan *types.Header)
	blocks, err := client.SubscribeNewHead(ctx, headers)
	if err != nil {
		cancel()
		return err
	}

	logs := make(chan types.Log)
	logger, err := client.SubscribeFilterLogs(ctx, ethereum.FilterQuery{Topics: topics}, logs)
	if err != nil {
		cancel()
		return err
	}

	for {
		select {
		case err := <-logger.Err():
			cancel()
			return err
		case err := <-blocks.Err():
			cancel()
			return err
		case header := <-headers:
			// get block details on new block header
			t := time.Unix(int64(header.Time), 0)
			ts := t.UTC().Format("2006-01-02T15:04:05Z07:00")
			ut := strconv.FormatInt(t.UTC().Unix(), 10)
			fmt.Println(header.Number, ts, ut)
		case log := <-logs:
			address := log.Address.String()
			fmt.Println(address)
		}
	}
}
