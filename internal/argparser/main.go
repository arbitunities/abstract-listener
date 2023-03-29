package argparser

import (
	"flag"
	"fmt"
	"os"
)

func ParseFlags() (*bool, *int64, *int64) {
	debug := flag.Bool("debug", false, "Debug log level ")
	from := flag.Int64("from", 0, "From block ")
	to := flag.Int64("to", 0, "To block. Default to latest")
	flag.Parse()

	if *from == 0 {
		fmt.Println("missing from block")
		os.Exit(1)
	}
	if *to == 0 {
		fmt.Println("missing to block")
		os.Exit(1)
	}
	return debug, from, to
}
