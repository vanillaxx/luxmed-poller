package luxmed

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	internalhttp "github.com/vanillaxx/luxmed-poller/internal/http"
	"github.com/vanillaxx/luxmed-poller/internal/luxmed/auth"
	"github.com/vanillaxx/luxmed-poller/internal/terms"
)

const (
	newReservationURL = "https://portalpacjenta.luxmed.pl/PatientPortal/NewPortal/terms/index"
)

type Client struct {
	*http.Client
	auth.User
}

func New(username, password string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("read request body: %w", err)
	}
	return &Client{Client: &http.Client{Jar: jar}, User: auth.User{Username: username, Password: password}}, nil
}

func getAvailableVisitsRequestParams(serviceVariantID, from, to string) map[string]string {
	return map[string]string{
		"cityId":           "3",
		"payerId":          "123",
		"serviceVariantId": serviceVariantID,
		"languageId":       "10",
		"searchDateFrom":   from,
		"searchDateTo":     to,
	}
}

// GetAvailableTerms returns available terms for given date range and servie.
func (c *Client) GetAvailableTerms(serviceVariantID, from, to string) (*terms.Info, error) {
	t, err := auth.GetToken(c.Client, c.User)
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}
	if err = auth.Login(c.Client, t); err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}

	data := url.Values{}
	params := getAvailableVisitsRequestParams(serviceVariantID, from, to)
	for k, v := range params {
		data.Set(k, v)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", newReservationURL, data.Encode()), nil)
	if err != nil {
		return nil, fmt.Errorf("new GET request: %W", err)
	}
	for h, v := range internalhttp.GetHeaders() {
		req.Header.Set(h, v)
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", t.TokenType, t.AccessToken))

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send http request: %w", err)
	}
	defer res.Body.Close()

	bytes, err := internalhttp.Decompress(res)
	if err != nil {
		return nil, fmt.Errorf("decompress: %w", err)
	}

	var result terms.Info
	if err = json.Unmarshal(bytes, &result); err != nil {
		return nil, fmt.Errorf("unmarshal visit terms: %s", err)
	}

	return &result, nil
}
