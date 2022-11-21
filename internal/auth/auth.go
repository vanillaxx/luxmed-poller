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
	LoginUrl       = "https://portalpacjenta.luxmed.pl/PatientPortalMobileAPI/api/token"
	ReservationURL = "https://portalpacjenta.luxmed.pl/PatientPortalMobileAPI/api/visits/available-terms"
)

func getHeaders() map[string]string {
	return map[string]string{
		"User-Agent":              "PatientPortal/4.14.0 (pl.luxmed.pp.LUX-MED; build:853; iOS 13.5.1) Alamofire/4.9.1",
		"Custom-User-Agent":       "PatientPortal; 4.14.0; 4380E6AC-D291-4895-8B1B-F774C318BD7D; iOS; 13.5.1; iPhone8,1",
		"Accept-Language":         "en;q=1.0, en-PL;q=0.9, pl-PL;q=0.8, ru-PL;q=0.7, uk-PL;q=0.6",
		"Accept-Encoding":         "gzip;q=1.0, compress;q=0.5",
		"Content-Type":            "application/x-www-form-urlencoded",
		"x-api-client-identifier": "iPhone",
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

func (c *LuxmedClient) AuthenticatedRequest(u User, method string, urlAddress string) (*http.Request, error) {
	t, err := c.GetToken(u)
	if err != nil {
		return nil, err
	}

	data := url.Values{}
	params := map[string]string{}
	params["cityId"] = "3"
	params["payerId"] = "123"
	params["serviceId"] = "4430"
	//params["LanguageId"] = "10"
	params["FromDate"] = "2022-11-18" // time.Now().Format(time.RFC3339Nano)
	params["ToDate"] = "2022-11-28"   // time.Now().AddDate(0, 0, 10).Format(time.RFC3339Nano)
	for k, v := range params {
		data.Set(k, v)
	}
	fmt.Printf("%s\n\n", data.Encode())

	req, err := http.NewRequest(method, fmt.Sprintf("%s?%s", urlAddress, data.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("err while creating new request: %s", err)
	}
	for h, v := range getHeaders() {
		req.Header.Set(h, v)
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", t.TokenType, t.AccessToken))

	return req, nil
}

// GetToken returns token for given username and
func (lc *LuxmedClient) GetToken(u User) (*oauth2.Token, error) {
	headers := getHeaders()
	params := getParams(u)

	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}

	req, err := http.NewRequest("POST", LoginUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Err while creating new request: %s", err)
	}
	for h, v := range headers {
		req.Header.Set(h, v)
	}

	// cj, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	lc.Client = &http.Client{
		// Jar: cj,
	}
	// lc := LuxmedClient{client: client}
	res, err := lc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Err while sending POST request: %s", err)
	}
	defer res.Body.Close()
	// res.Cookies()

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
