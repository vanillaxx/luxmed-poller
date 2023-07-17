package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	internalhttp "github.com/vanillaxx/luxmed-poller/internal/http"
	"golang.org/x/oauth2"
)

const (
	TokenUrl = "https://portalpacjenta.luxmed.pl/PatientPortalMobileAPI/api/token"
	LoginUrl = "https://portalpacjenta.luxmed.pl/PatientPortal/Account/LogInToApp"
)

type User struct {
	Username string
	Password string
}

func getTokenRequestParams(u User) map[string]string {
	return map[string]string{
		"client_id":  "iPhone",
		"grant_type": "password",
		"username":   u.Username,
		"password":   u.Password,
	}
}

func getLoginRequestParams() map[string]string {
	return map[string]string{
		"app":              "search",
		"client":           "3",
		"paymentSupported": "true",
		"lang":             "pl",
	}
}

// GetToken returns token for given user.
func GetToken(c *http.Client, u User) (*oauth2.Token, error) {
	headers := internalhttp.GetHeaders()
	params := getTokenRequestParams(u)

	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}

	req, err := http.NewRequest("POST", TokenUrl, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("new POST request: %w", err)
	}
	for h, v := range headers {
		req.Header.Set(h, v)
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send HTTP request: %w", err)
	}
	defer res.Body.Close()

	bytes, err := internalhttp.Decompress(res)
	if err != nil {
		return nil, fmt.Errorf("decompress: %w", err)
	}

	t := &oauth2.Token{}
	if err = json.Unmarshal(bytes, &t); err != nil {
		return nil, fmt.Errorf("unmarshall token: %w", err)
	}

	return t, nil
}

// Login logs in to luxmed portal.
func Login(c *http.Client, t *oauth2.Token) error {
	data := url.Values{}
	params := getLoginRequestParams()
	for k, v := range params {
		data.Set(k, v)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", LoginUrl, data.Encode()), nil)
	if err != nil {
		return fmt.Errorf("new GET request: %w", err)
	}

	for h, v := range internalhttp.GetHeaders() {
		req.Header.Set(h, v)
	}
	req.Header.Set("Authorization", t.AccessToken)

	res, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("send HTTP request: %s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code %d not expected during login attempt", res.StatusCode)
	}
	return nil
}
