package erc4337

import (
	"strings"

	"github.com/arbitunities/abstract-listener/pkg/config"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func ABI() *abi.ABI {
	it, err := abi.JSON(strings.NewReader(Erc4337MetaData.ABI))
	if err != nil {
		return &abi.ABI{}
	}
	return &it
}

func NewErc4337Singleton(backend bind.ContractBackend) (*Erc4337, error) {
	contract, err := bindErc4337(config.EntryPointAddress, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erc4337{Erc4337Caller: Erc4337Caller{contract: contract}, Erc4337Transactor: Erc4337Transactor{contract: contract}, Erc4337Filterer: Erc4337Filterer{contract: contract}}, nil
}
