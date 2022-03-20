package main

import (
	"fmt"
	"os"

	"coinconv/service/coinmarketcap"
)

func main() {
	args := os.Args
	if len(args) != 4 {
		exit(1, "use format <value> <from currency> <to currency>")
	}

	m, err := coinmarketcap.GetCoinMarketCap()
	if err != nil {
		exit(2, err.Error())
	}

	curr1, ok := m[args[2]]
	if !ok {
		exit(3, "can't find currency "+args[2])
	}

	curr2, ok := m[args[3]]
	if !ok {
		exit(4, "can't find currency "+args[3])
	}

	val, err := coinmarketcap.PriceConversion(args[1], curr1, curr2)
	if err != nil {
		exit(1, err.Error())
	}

	fmt.Print(val)
}

func exit(code int, message string) {
	fmt.Fprintf(os.Stderr, message)
	os.Exit(code)
}
