package http

import (
	"github.com/lamg/tesis"
	a "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestHTTPPortal(t *testing.T) {
	var e error
	var fs *ServFS
	// { files referenced in fs exist
	// in cwd }
	fs = &ServFS{"cert.pem", "key.pem"}
	// {fs initialized}
	hp := "localhost:10443"
	au := &tesis.DummyAuth{}
	qr := tesis.NewDummyManager()
	go ListenAndServe(hp, au, qr, fs)
	time.Sleep(1 * time.Second)

	var j bool
	var cl *PortalUser

	cl = NewPortalUser(hp)
	j, e = cl.Auth("a", "a")
	a.NoError(t, e, "Auth failed")
	if a.NotNil(t, j) {
		t.Logf("Valid: %t", j)
	}

	var s string
	if e == nil {
		s, e = cl.Sync()
	}
	if a.NoError(t, e) {
		t.Log(s)
	}
	if e == nil {
		s, e = cl.Sync()
		// s has non nil error field
		// since cl.Sync() tries to
		// synchronize the same record again
	}
	if a.NoError(t, e) {
		t.Log(s)
	}
}
