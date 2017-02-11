package tesis

import (
	a "github.com/stretchr/testify/assert"

	"testing"
)

/*var (
	c   *http.Client
	url = "https://localhost:10443"
)*/

func TestHTTPPortal(t *testing.T) {
	hp := "localhost:10443"
	ce := "cert.pem"
	ke := "key.pem"
	rs := []Route{
		Route{"Root", "GET", "/", rootH},
		Route{"Auth", "POST", "/", authH},
		Route{"Conv", "GET", "/conv", convH}}
	
	s, e := NewHTTPPortal(hp, ce, ke, rs)
	a.NoError(t, e, "Error creating server")
	a.HTTPError(t, rootH, "GET", "", nil)

	cr := &Credentials{user: user, pass: password}
	j, e := s.Auth(cr)
	a.NoError(t, e, "Auth failed")
	t.Logf("raw: %s", j.Raw)
}
