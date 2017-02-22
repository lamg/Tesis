package tesis

import (
	"crypto/rand"
	"crypto/rsa"
	a "github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type DummyAuth struct {
}

func (d *DummyAuth) Authenticate(u, p string) (b bool) {
	b = u == p
	return
}

type DummyManager struct {
}

func (m *DummyManager) UserInfo(u string) (inf *Info, e error) {
	inf = &Info{
		SentMessages: 18,
		RecvMessages: 40,
		MailStorage:  67,
		InternetDwnl: 87,
		WifiLogons: []WifiL{
			WifiL{
				Ip:    "192.168.0.10",
				Place: "Rectoría",
				Date:  time.Now(),
			},
		},
	}
	return
}

func TestHTTPPortal(t *testing.T) {
	var e error
	hp := "localhost:10443"
	ce := "serv/cert.pem"
	ke := "serv/key.pem"
	au := &DummyAuth{}
	qr := &DummyManager{}
	_, e = NewHTTPPortal(hp, ce, ke, au, qr)
	a.NoError(t, e, "Error creating server")
	a.HTTPError(t, rootH, "GET", "", nil)

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
}

func TestGenKey(t *testing.T) {
	k, e := rsa.GenerateKey(rand.Reader, 2048)
	a.True(t, k != nil && e == nil)
}
