package config

import (
	"github.com/ethereum/go-ethereum/common"
)

const (
	EntryPointAddressString = "0x0576a174D229E3cFA37253523E645A78A0C91B57"
)

var (
	EntryPointAddress = common.HexToAddress(EntryPointAddressString)
)
