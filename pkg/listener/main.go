package listener

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/arbitunities/abstract-listener/pkg/erc4337"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Listener struct {
	Client *ethclient.Client
}

func NewListener(client *ethclient.Client) *Listener {
	return &Listener{
		Client: client,
	}
}

func (l *Listener) FilterLogs(ctx context.Context, filter ethereum.FilterQuery) ([]types.Log, error) {
	fmt.Println("Filtering logs...")
	return l.Client.FilterLogs(ctx, filter)
}

func (l *Listener) GetAddressTransactions(
	ctx context.Context,
	filter *common.Address,
	fromBlock, toBlock int,
) (txs []*ParsedTransaction, err error) {
	// txs := []*types.Receipt{}
	for i := fromBlock; i <= int(toBlock); i++ {
		fmt.Println(i)
		block, err := l.Client.BlockByNumber(ctx, big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}
		for _, tx := range block.Transactions() {
			if tx.To().String() == filter.String() {
				txs = append(txs, NewParsedTransaction(tx))
				receipt := l.TransactionReceipt(ctx, tx.Hash())
				DecodeTransactionLogs(receipt, erc4337.ABI())
			}
		}
	}
	return txs, nil
}

type ParsedTransaction struct {
	Hash     string
	ChainId  *big.Int
	Value    string
	From     string
	To       string
	Gas      uint64
	GasPrice uint64
	Nonce    uint64
	Data     string
}

func GetTransactionSender(tx *types.Transaction) common.Address {
	from, err := types.Sender(types.NewEIP155Signer(tx.ChainId()), tx)
	if err != nil {
		from, err := types.Sender(types.HomesteadSigner{}, tx)
		if err != nil {
			return common.HexToAddress("0x404")
		}
		return from
	}
	return from
}

func DecodeTransactionLogs(receipt *types.Receipt, contractABI *abi.ABI) {
	for _, vLog := range receipt.Logs {
		// topic[0] is the event name
		event, err := contractABI.EventByID(vLog.Topics[0])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Event Name: %s\n", event.Name)
		// topic[1:] is other indexed params in event
		if len(vLog.Topics) > 1 {
			for i, param := range vLog.Topics[1:] {
				fmt.Printf("Indexed params %d in hex: %s\n", i, param)
				fmt.Printf("Indexed params %d decoded %s\n", i, common.HexToAddress(param.Hex()))
			}
		}

		if len(vLog.Data) > 0 {
			fmt.Printf("Log Data in Hex: %s\n", hex.EncodeToString(vLog.Data))
			outputDataMap := make(map[string]interface{})
			err = contractABI.UnpackIntoMap(outputDataMap, event.Name, vLog.Data)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Event outputs: %v\n", outputDataMap)
		}

	}
}

func NewParsedTransaction(tx *types.Transaction) *ParsedTransaction {
	return &ParsedTransaction{
		Hash:     tx.Hash().Hex(),
		ChainId:  tx.ChainId(),
		Value:    tx.Value().String(),
		To:       tx.To().Hex(),
		From:     GetTransactionSender(tx).String(),
		Nonce:    tx.Nonce(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice().Uint64(),
		Data:     hex.EncodeToString(tx.Data()),
	}
}

func (l *Listener) TransactionReceipt(ctx context.Context, hash common.Hash) *types.Receipt {
	receipt, err := l.Client.TransactionReceipt(ctx, hash)
	if err != nil {
		return nil
	}
	return receipt
}

func (l *Listener) SubscribeNewHead(ctx context.Context) error {
	headers := make(chan *types.Header)
	blocks, err := l.Client.SubscribeNewHead(ctx, headers)
	if err != nil {
		return err
	}

	for {
		select {
		case err := <-blocks.Err():
			return err
		case header := <-headers:
			// get block details on new block header
			t := time.Unix(int64(header.Time), 0)
			ts := t.UTC().Format("2006-01-02T15:04:05Z07:00")
			ut := strconv.FormatInt(t.UTC().Unix(), 10)
			fmt.Println(header.Number, ts, ut)
		}
	}
}

func (l *Listener) SubscribeFilterLogs(filter ethereum.FilterQuery) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	logs := make(chan types.Log)
	logger, err := l.Client.SubscribeFilterLogs(ctx, filter, logs)
	if err != nil {
		cancel()
		return err
	}

	for {
		select {
		case err := <-logger.Err():
			cancel()
			return err
		case log := <-logs:
			address := log.Address.String()
			fmt.Println(address)
		}
	}
}
