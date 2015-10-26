// Package httpfetch is an illustrative package used for
// a talk on testing in Go. The talk is at http://godoc.org/github.com/rogpeppe/talks/testing.talk/testing.slide .
package httpfetch

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var httpGet = http.Get

// GetURLAsString makes a GET request to the
// given URL and returns the result as a string.
func GetURLAsString(url string) (string, error) {
	resp, err := httpGet(url)
	if err != nil {
		return "", fmt.Errorf("GET failed: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GET returned unexpected status %q", resp.Status)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read body: %v", err)
	}
	return string(data), nil
}
