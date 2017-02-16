package tesis

import (
	"crypto/rand"
	"crypto/rsa"
	a "github.com/stretchr/testify/assert"
	"testing"
)

type DummyAuth struct {
}

func (d *DummyAuth) Authenticate(u, p string) (b bool) {
	b = u == p
	return
}

func TestHTTPPortal(t *testing.T) {
	hp := "localhost:10443"
	ce := "cert.pem"
	ke := "key.pem"
	au := &DummyAuth{}
	_, e := NewHTTPPortal(hp, ce, ke, au)
	a.NoError(t, e, "Error creating server")
	a.HTTPError(t, rootH, "GET", "", nil)

	cr := &Credentials{User: "a", Pass: "a"}
	cl := NewPortalUser(hp)
	j, e := cl.Auth(cr)
	a.NoError(t, e, "Auth failed")
	if a.NotNil(t, j) {
		t.Logf("Valid: %t", j)
	}
	j, e = cl.Check()
	if !a.True(t, j && e == nil) {
		t.Logf("%t %s", j, e.Error())
	}
}

func TestGenKey(t *testing.T) {
	k, e := rsa.GenerateKey(rand.Reader, 2048)
	a.True(t, k != nil && e == nil)
}
