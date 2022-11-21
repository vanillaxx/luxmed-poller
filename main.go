package main

import (
	"fmt"
	"os"
	"time"

	"github.com/vanillaxx/luxmed-poller/internal/client"
)

func main() {
	u := client.User{
		Username: "username",
		Password: "password",
	}
	lc := client.LuxmedClient{}
	data := &client.Params{
		"cityId":    "3",
		"payerId":   "123",
		"serviceId": "4430",
		"FromDate":  time.Now().Format(time.RFC3339Nano),
		"ToDate":    time.Now().AddDate(0, 0, 10).Format(time.RFC3339Nano),
	}
	vt, err := lc.GetVisitTerms(u, data)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	fmt.Printf("%+v", vt)
}
