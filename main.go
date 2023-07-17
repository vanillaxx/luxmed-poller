package main

import (
	"fmt"
	"os"

	"github.com/vanillaxx/luxmed-poller/internal/luxmed"
)

func main() {
	u, p := "username", "password"
	lc, err := luxmed.New(u, p)
	if err != nil {
		fmt.Printf("new: %s", err)
		os.Exit(1)
	}
	serviceVariantId := "4502"
	from := "2023-07-18"
	to := "2023-07-30"
	result, err := lc.GetAvailableTerms(serviceVariantId, from, to)
	if err != nil {
		fmt.Printf("get available visits: %s", err)
		os.Exit(1)
	}

	fmt.Printf("result: %v", result)
}
