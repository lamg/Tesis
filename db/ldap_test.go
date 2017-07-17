package db

import (
	"github.com/lamg/tesis"
	a "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// user
// password
// account suffix
// ldap server address
var u, p, lda, sf string

func init() {
	u, p, lda, sf = os.Getenv("AD_USER"), os.Getenv("AD_PASS"),
		os.Getenv("AD_ADDR"), os.Getenv("AD_SUFF")
}

func TestLDAPAuth(t *testing.T) {
	var l *LDAPAuth
	var e error

	l, e = NewLDAPAuth(lda, sf)
	if a.NoError(t, e) {
		var b bool
		b, e = l.Authenticate(u, p)
		a.NoError(t, e)
		a.True(t, b, "Failed authentication")
	}
}

func TestGetUsers(t *testing.T) {
	var us []tesis.DBRecord
	var r tesis.RecordProvider
	var e error
	r, e = NewLDAPProv(u, p, lda, -1)
	if a.NoError(t, e) {
		us, e = r.Records()
	}
	if a.NoError(t, e) {
		t.Log(len(us))
		for _, j := range us {
			if j.Name == "Luis Angel Mendez Gort" {
				t.Log(j)
			}
		}
	}
}

//go test -v -run 'TestLDAPAuth|TestGetUsers'

func TestGetLDAPEntry(t *testing.T) {
	var l *LDAPAuth
	var e error

	l, e = NewLDAPAuth(lda, sf)
	var b bool
	if a.NoError(t, e) {
		b, e = l.Authenticate(u, p)
	}
	var r []string
	if a.True(t, b, "Failed authentication") {
		r, e = Search("Luis Angel Mendez Gort", l.c)
	}
	if a.NoError(t, e) {
		for _, j := range r {
			t.Log(j)
		}
	}
}

func TestUserInfo(t *testing.T) {
	var l *LDAPAuth
	var e error
	l, e = NewLDAPAuth(lda, sf)
	var b bool
	if a.NoError(t, e) {
		b, e = l.Authenticate(u, p)
	}
	var ui *tesis.UserInfo
	if a.True(t, b) {
		ui, e = l.UserInfo(u)
	}
	a.True(t, e == nil && ui.Name != "" &&
		ui.UserName == u)
}

func TestUpdate(t *testing.T) {
	var r *LDAPAuth
	var e error
	r, e = NewLDAPAuth(lda, sf)
	var b bool
	if a.NoError(t, e) {
		b, e = r.Authenticate(u, p)
	}
	var dr *tesis.DBRecord
	if a.True(t, e == nil && b) {
		dr = &tesis.DBRecord{
			Name: "Luis Ángel Méndez Gort",
			IN:   "89011914982",
			Addr: "Briones Montoto",
			Tel:  "48791438",
		}
		e = r.Update("CN=Luis Angel Mendez Gort,OU=_GrupoRedes,DC=upr,DC=edu,DC=cu", dr)
	}
	var rc *tesis.DBRecord
	if a.NoError(t, e) {
		rc, e = r.UserRecord("luis.mendez")
	}
	if a.NoError(t, e) {
		a.True(t, dr.IN == rc.IN && dr.Addr == rc.Addr &&
			dr.Tel == rc.Tel)
	}
}
