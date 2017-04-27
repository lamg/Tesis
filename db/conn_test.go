package db

import (
	"github.com/lamg/tesis"
	a "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	//ldap server address
	lds = "ad.upr.edu.cu"
	//account suffix
	sf = "@upr.edu.cu"
	//ldap server port
	ldp = 636
)

var u, p string
var l *LDAPAuth
var e error

func init() {
	u, p = os.Getenv("UPR_USER"), os.Getenv("UPR_PASS")
}

func TestLDAPAuth(t *testing.T) {
	l, e = NewLDAPAuth(lds, sf, ldp)
	if a.NoError(t, e) {
		var b bool
		b, e = l.Authenticate(u, p)
		a.True(t, b, "Failed authentication")
	}
}

func TestGetUsers(t *testing.T) {
	var us []tesis.DBRecord
	if a.NoError(t, e) {
		us, e = l.GetUsers()
	}
	if a.NoError(t, e) {
		t.Log(len(us))
	}
}

//go test -v -run 'TestLDAPAuth|TestGetUsers'
