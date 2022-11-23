package client

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/vanillaxx/luxmed-poller/internal/visit"
	"golang.org/x/oauth2"
)

const (
	BaseUrl       = "https://portalpacjenta.luxmed.pl/PatientPortalMobileAPI/api"
	TokenEndpoint = "token"
	TermsEndpoint = "visits/available-terms"
)

type User struct {
	Username string
	Password string
}

type LuxmedClient struct {
	*http.Client
}

// NewLuxmedClient returns client to luxmed API
func NewLuxmedClient() LuxmedClient {
	return LuxmedClient{&http.Client{}}
}

type Params map[string]string

func (p *Params) mapToUrlParams() string {
	data := url.Values{}
	for k, v := range *p {
		data.Set(k, v)
	}
	return data.Encode()
}

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

func getParams(u User) Params {
	return map[string]string{
		"client_id":  "iPhone",
		"grant_type": "password",
		"username":   u.Username,
		"password":   u.Password,
	}
}

// GetVisitTerms returns a list of visit terms for given parameters
func (lc *LuxmedClient) GetVisitTerms(u User, p *Params) ([]visit.VisitTerm, error) {
	req, err := lc.authenticatedRequest(u, http.MethodGet, fmt.Sprintf("%s/%s", BaseUrl, TermsEndpoint), p)
	if err != nil {
		return nil, fmt.Errorf("create request for getting visit terms: %s", err)
	}
	res, err := lc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute request for getting visit terms: %s", err)
	}
	defer res.Body.Close()
	gzreader, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("create gzip reader: %s", err)
	}
	defer gzreader.Close()
	bytes, err := ioutil.ReadAll(gzreader)
	if err != nil {
		return nil, fmt.Errorf("read visit terms with gzip reader: %s", err)
	}
	if err != nil {
		return nil, err
	}
	resp := &visit.VisitTermsResponse{}
	if err = json.Unmarshal(bytes, resp); err != nil {
		return nil, fmt.Errorf("unmarshal visit terms: %s", err)
	}
	return resp.VisitTerms, nil
}

// authenticatedRequest creates an Oauth2 authenticated request
func (lc *LuxmedClient) authenticatedRequest(u User, method, url string, p *Params) (*http.Request, error) {
	t, err := lc.getToken(u)
	if err != nil {
		return nil, err
	}

	if len(*p) > 0 {
		url = fmt.Sprintf("%s?%s", url, p.mapToUrlParams())
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("err while creating new request: %s", err)
	}
	for h, v := range getHeaders() {
		req.Header.Set(h, v)
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", t.TokenType, t.AccessToken))
	return req, nil
}

// getToken returns authentication token for given User
func (lc *LuxmedClient) getToken(u User) (*oauth2.Token, error) {
	headers := getHeaders()
	p := getParams(u)

	encodedParams := p.mapToUrlParams()
	wholeUrl := fmt.Sprintf("%s/%s", BaseUrl, TokenEndpoint)
	req, err := http.NewRequest(http.MethodPost, wholeUrl, strings.NewReader(encodedParams))
	if err != nil {
		return nil, fmt.Errorf("err while creating new request: %s", err)
	}
	for h, v := range headers {
		req.Header.Set(h, v)
	}

	res, err := lc.Do(req)
	if err != nil {
		return nil, fmt.Errorf("err while sending POST request: %s", err)
	}
	defer res.Body.Close()

	gzreader, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("err while creating gzip reader: %s", err)
	}
	defer gzreader.Close()
	bytes, err := ioutil.ReadAll(gzreader)
	if err != nil {
		return nil, fmt.Errorf("err while reading with gzip reader: %s", err)
	}
	resp := &oauth2.Token{}
	if err = json.Unmarshal(bytes, &resp); err != nil {
		return nil, fmt.Errorf("err while unmarshalling token: %s", err)
	}
	return resp, nil
}
