package http

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
)

func GetHeaders() map[string]string {
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

func Decompress(res *http.Response) ([]byte, error) {
	gzreader, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("new gzip reader: %w", err)
	}
	defer gzreader.Close()
	bytes, err := io.ReadAll(gzreader)
	if err != nil {
		return nil, fmt.Errorf("read all: %w", err)
	}
	return bytes, nil
}
