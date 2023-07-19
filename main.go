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
	services, err := lc.GetAvailableServices()
	if err != nil {
		fmt.Printf("get available services: %s", err)
		os.Exit(1)
	}
	fmt.Printf("available services: %v", services)

	serviceVariantId := "4565"
	from := "2023-07-19"
	to := "2023-07-30"
	terms, err := lc.GetAvailableTerms(serviceVariantId, from, to)
	if err != nil {
		fmt.Printf("get available visits: %s", err)
		os.Exit(1)
	}

	fmt.Printf("result: %v", terms)
}
