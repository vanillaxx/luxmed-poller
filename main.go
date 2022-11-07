package main

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/oauth2"
)

func main() {
	u, p := "username", "password"
	t, err := getToken(u, p)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}
	fmt.Printf("%s", t)
}

func getToken(username string, password string) (oauth2.Token, error) {
	headers := map[string]string{
		"User-Agent":              "PatientPortal/4.14.0 (pl.luxmed.pp.LUX-MED; build:853; iOS 13.5.1) Alamofire/4.9.1",
		"Custom-User-Agent":       "PatientPortal; 4.14.0; 4380E6AC-D291-4895-8B1B-F774C318BD7D; iOS; 13.5.1; iPhone8,1",
		"Accept-Language":         "en;q=1.0, en-PL;q=0.9, pl-PL;q=0.8, ru-PL;q=0.7, uk-PL;q=0.6",
		"Accept-Encoding":         "gzip;q=1.0, compress;q=0.5",
		"Content-Type":            "application/x-www-form-urlencoded",
		"x-api-client-identifier": "iPhone",
	}
	params := map[string]string{
		"client_id":  "iPhone",
		"grant_type": "password",
		"username":   username,
		"password":   password,
	}
	baseUrl := "https://portalpacjenta.luxmed.pl/PatientPortalMobileAPI/api/token"

	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}

	req, err := http.NewRequest("POST", baseUrl, strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Printf("Err while creating new request: %s", err)
		os.Exit(1)
	}
	for h, v := range headers {
		req.Header.Set(h, v)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Err while sending POST request: %s", err)
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
		fmt.Println(err)
		os.Exit(1)
	}
	resp := oauth2.Token{}
	if err = json.Unmarshal(bytes, &resp); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return resp, nil
}
