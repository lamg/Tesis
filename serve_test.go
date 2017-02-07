package tesis

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	a "github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

var (
	c   *http.Client
	url = "https://localhost:10443"
)

func init() {
	go serveTLS() //what if it fails to start
	//start server

	cfg := &tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{TLSClientConfig: cfg}
	c = &http.Client{Transport: tr}
	time.Sleep(1000000 * time.Nanosecond)
}

func TestServe(t *testing.T) {
	r, e := c.Get(url)
	//client make request
	if a.NoError(t, e) {
		bd, e := ioutil.ReadAll(r.Body)
		if a.NoError(t, e) {
			s := string(bd)
			t.Logf("s: %s", s)
		}
	}
	//analyze response
	//close server
}

func TestAuth(t *testing.T) {
	cr := Credentials{user: user, pass: password}
	b, e := json.Marshal(cr)
	if a.NoError(t, e) {
		br := bytes.NewReader(b)
		r, e := c.Post(url, "application/json", br)
		//use assert.HTTPError
		if a.NoError(t, e) {
			bd, e := ioutil.ReadAll(r.Body)
			if a.NoError(t, e) {
				s := string(bd)
				t.Logf("s: %s", s)
			}
		}
	}
}
