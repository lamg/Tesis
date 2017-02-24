package tesis

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	a "github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestHTTPPortal(t *testing.T) {
	var e error
	var h *HTTPPortal
	os.Chdir("serv")
	// { files referenced in http_serve.go exist
	// in cwd }
	hp := "localhost:10443"
	au := &DummyAuth{}
	qr := &DummyManager{}
	h, e = NewHTTPPortal(hp, au, qr)
	a.NoError(t, e, "Error creating server")
	go h.Serve()
	time.Sleep(1 * time.Second)

	var j bool
	var cr *Credentials
	var cl *PortalUser
	cr = &Credentials{User: "a", Pass: "a"}
	cl = NewPortalUser(hp)
	j, e = cl.Auth(cr)
	a.NoError(t, e, "Auth failed")
	if a.NotNil(t, j) {
		t.Logf("Valid: %t", j)
	}
	var inf *Info
	inf, e = cl.Info()
	if a.NoError(t, e) {
		t.Logf("%v", inf)
	}

	var s string
	s, e = cl.Index()
	if a.NoError(t, e) {
		t.Log(s)
	}

	var c context.Context
	c = context.Background()
	h.Shutdown(c)
}

func TestGenKey(t *testing.T) {
	k, e := rsa.GenerateKey(rand.Reader, 2048)
	a.True(t, k != nil && e == nil)
}
