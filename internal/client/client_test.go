package client

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
)

func TestGetToken(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/token", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		body, _ := json.Marshal(&oauth2.Token{TokenType: "Bearer", AccessToken: "GoGoPowerRangers"})
		gw := gzip.NewWriter(w)
		gw.Write([]byte(body))
		defer gw.Close()
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	url := ts.URL
	c := NewLuxmedClient(url)
	u := User{"foo", "bar"}
	token, err := c.getToken(u)
	assert.Equal(t, "Bearer", token.TokenType)
	assert.Equal(t, "GoGoPowerRangers", token.AccessToken)
	require.NoError(t, err)
}

func TestAuthenticatedRequest(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		body, _ := json.Marshal(&oauth2.Token{TokenType: "Bearer", AccessToken: "GoGoPowerRangers"})
		gw := gzip.NewWriter(w)
		gw.Write([]byte(body))
		defer gw.Close()
	})
	ts := httptest.NewServer(mux)
	defer ts.Close()
	url := ts.URL
	c := NewLuxmedClient(url)
	u := User{"foo", "bar"}
	req, err := c.authenticatedRequest(u, http.MethodGet, "https://test.pl/test", &Params{"param1": "value1"})
	assert.Contains(t, req.Header, "Authorization")
	assert.Equal(t, []string{"Bearer GoGoPowerRangers"}, req.Header["Authorization"])
	require.NoError(t, err)
}
