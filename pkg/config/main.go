package config

import (
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
)

type Config struct {
	Rpc                  string
	EntryPointAddress    common.Address
	SimpleAccountFactory common.Address
}

func NewConfig() *Config {
	godotenv.Load(".env")
	return &Config{
		Rpc:                  os.Getenv("RPC"),
		EntryPointAddress:    common.HexToAddress("0x0576a174D229E3cFA37253523E645A78A0C91B57"),
		SimpleAccountFactory: common.HexToAddress("0x71D63edCdA95C61D6235552b5Bc74E32d8e2527B"),
	}
}
