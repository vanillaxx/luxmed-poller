package main

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/vanillaxx/luxmed-poller/internal/auth"
)

func main() {
	u := auth.User{
		Username: "username",
		Password: "password",
	}
	lc := auth.LuxmedClient{}
	req, err := lc.AuthenticatedRequest(u, "GET", auth.ReservationURL)
	if err != nil {
		fmt.Printf("Err while sending GET request: %s", err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n\n", req)
	res, err := lc.Do(req)
	if err != nil {
		fmt.Printf("Err while sending GET request: %s", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	gzreader, err := gzip.NewReader(res.Body)
	if err != nil {
		fmt.Printf("Err while creating gzip reader: %s", err)
		os.Exit(1)
	}
	defer gzreader.Close()
	bytes, err := ioutil.ReadAll(gzreader)
	if err != nil {
		fmt.Printf("Err while reading with gzip reader: %s", err)
		os.Exit(1)
	}
	sss := string(bytes)

	//t, err := auth.GetToken(u)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	fmt.Printf("%+v", sss)
}
