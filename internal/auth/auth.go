package auth

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

type User struct {
	Username string
	Password string
}

type LuxmedClient struct {
	*http.Client
}

const (
	TokenUrl          = "https://portalpacjenta.luxmed.pl/PatientPortalMobileAPI/api/token"
	LoginUrl          = "https://portalpacjenta.luxmed.pl/PatientPortal/Account/LogInToApp"
	ReservationURL    = "https://portalpacjenta.luxmed.pl/PatientPortalMobileAPI/api/visits/available-terms"
	NewReservationURL = "https://portalpacjenta.luxmed.pl/PatientPortal/NewPortal/terms/index"
)

func getHeaders() map[string]string {
	return map[string]string{
		"User-Agent":              "okhttp/3.11.0",
		"Custom-User-Agent":       "PatientPortal; 4.19.0; 4380E6AC-D291-4895-8B1B-F774C318BD7D; iOS; 13.5.1; iPhone8,1",
		"Accept":                  "application/json, text/plain, */*",
		"Accept-Language":         "en;q=1.0, en-PL;q=0.9, pl-PL;q=0.8, ru-PL;q=0.7, uk-PL;q=0.6",
		"Accept-Encoding":         "gzip;q=1.0, compress;q=0.5",
		"Content-Type":            "application/x-www-form-urlencoded",
		"x-api-client-identifier": "iPhone",
		"Host":                    "portalpacjenta.luxmed.pl",
		"Origin":                  "https://portalpacjenta.luxmed.pl",
	}
}

func getParams(u User) map[string]string {
	return map[string]string{
		"client_id":  "iPhone",
		"grant_type": "password",
		"username":   u.Username,
		"password":   u.Password,
	}
}

// GetToken returns token for given username and
func (lc *LuxmedClient) GetToken(u User) (*oauth2.Token, error) {
	headers := getHeaders()
	params := getParams(u)

	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}

	req, err := http.NewRequest("POST", TokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Err while creating new request: %s", err)
	}
	for h, v := range headers {
		req.Header.Set(h, v)
	}

	res, err := lc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Err while sending POST request: %s", err)
	}
	defer res.Body.Close()

	gzreader, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Err while creating gzip reader: %s", err)
	}
	defer gzreader.Close()
	bytes, err := ioutil.ReadAll(gzreader)
	if err != nil {
		return nil, fmt.Errorf("Err while reading with gzip reader: %s", err)
	}
	resp := &oauth2.Token{}
	if err = json.Unmarshal(bytes, &resp); err != nil {
		return nil, fmt.Errorf("Err while unmarshalling token: %s", err)
	}

	return resp, nil
}

func (lc *LuxmedClient) Login(t *oauth2.Token) error {
	data := url.Values{}
	params := map[string]string{
		"app":              "search",
		"client":           "3",
		"paymentSupported": "true",
		"lang":             "pl",
	}
	for k, v := range params {
		data.Set(k, v)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", LoginUrl, data.Encode()), nil)
	if err != nil {
		return fmt.Errorf("err while creating new request: %s", err)
	}

	for h, v := range getHeaders() {
		req.Header.Set(h, v)
	}
	req.Header.Set("Authorization", t.AccessToken)

	res, err := lc.Do(req)
	if err != nil {
		return fmt.Errorf("err while sending request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("got status code: %d", res.StatusCode)
	}
	return nil
}

func (lc *LuxmedClient) GetAvailableVisits(u User) (interface{}, error) {
	t, err := lc.GetToken(u)
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}
	if err = lc.Login(t); err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	data := url.Values{}
	params := map[string]string{}
	params["cityId"] = "3"
	params["payerId"] = "123"
	params["serviceVariantId"] = "4502"
	params["languageId"] = "10"
	params["searchDateFrom"] = "2023-07-17" // time.Now().Format(time.RFC3339Nano)
	params["searchDateTo"] = "2023-07-30"   // time.Now().AddDate(0, 0, 10).Format(time.RFC3339Nano)
	for k, v := range params {
		data.Set(k, v)
	}
	fmt.Printf("%s\n\n", data.Encode())

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", NewReservationURL, data.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("err while creating new request: %s", err)
	}
	for h, v := range getHeaders() {
		req.Header.Set(h, v)
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", t.TokenType, t.AccessToken))

	return lc.readRequestBody(req)
}

func (lc *LuxmedClient) readRequestBody(req *http.Request) (string, error) {
	res, err := lc.Do(req)
	if err != nil {
		return "", fmt.Errorf("err while sending request: %s", err)
	}
	defer res.Body.Close()

	gzreader, err := gzip.NewReader(res.Body)
	if err != nil {
		return "", fmt.Errorf("err while creating gzip reader: %s", err)
	}
	defer gzreader.Close()
	bytes, err := ioutil.ReadAll(gzreader)
	if err != nil {
		return "", fmt.Errorf("err while reading with gzip reader: %s", err)
	}

	fmt.Printf("bytes: %s\n", bytes)

	var result string
	if err = json.Unmarshal(bytes, &result); err != nil {
		return "", fmt.Errorf("unmarshal visit terms: %s", err)
	}

	return result, nil
}
