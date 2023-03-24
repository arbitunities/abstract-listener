package chains

import "math/big"

// chains is a mapping of id to friendly name.
var chains = map[string]string{
	"1":   "eth",
	"5":   "goerli",
	"137": "polygon",
}

// NameById returns the chain name for the given id.
// Otherwise returns the chain as a string.
func NameById(id *big.Int) (name string) {
	if name, found := chains[id.String()]; found {
		return name
	}
	return id.String()
}
