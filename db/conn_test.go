package db

import (
	"github.com/stretchr/testify/assert"
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

func TestLDAPAuth(t *testing.T) {
	var e error
	var l *LDAPAuth
	u, p := os.Getenv("UPR_USER"), os.Getenv("UPR_PASS")
	l, e = NewLDAPAuth(lds, sf, ldp)
	if assert.NoError(t, e) {
		var b bool
		b, e = l.Authenticate(u, p)
		assert.True(t, b, "Failed authentication")
	}
}
