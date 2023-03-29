package erc4337

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

func ABI() *abi.ABI {
	it, err := abi.JSON(strings.NewReader(Erc4337ABI))
	if err != nil {
		return &abi.ABI{}
	}
	return &it
}
