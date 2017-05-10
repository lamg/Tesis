package db

import (
	"github.com/go-ldap/ldap"
	"github.com/lamg/tesis"
	a "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	//ldap server address
	lda = "ad.upr.edu.cu:636"
	//account suffix
	sf = "@upr.edu.cu"
)

var u, p string

func init() {
	u, p = os.Getenv("UPR_USER"), os.Getenv("UPR_PASS")
}

func TestLDAPAuth(t *testing.T) {
	var l *LDAPAuth
	var e error

	l, e = NewLDAPAuth(lda, sf)
	if a.NoError(t, e) {
		var b bool
		b, e = l.Authenticate(u, p)
		a.True(t, b, "Failed authentication")
	}
}

func TestGetUsers(t *testing.T) {
	var us []tesis.DBRecord
	var r tesis.RecordProvider
	var e error
	r, e = NewLDAPProv(u, p, "10.2.24.35:636", -1)
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
	var r *ldap.Entry
	if a.True(t, b, "Failed authentication") {
		r, e = Search("Luis Angel Mendez Gort", l.c)
	}
	if a.NoError(t, e) {
		for i := range r.Attributes {
			t.Logf("%s: %v", r.Attributes[i].Name, r.Attributes[i].Values)
		}
	}

}
