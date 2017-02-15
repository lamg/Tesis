package tesis

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
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
	u, p := os.Getenv("UPR_USER"), os.Getenv("UPR_PASS")
	l, e := NewLDAPAuth(lds, sf, ldp)
	if assert.NoError(t, e) {
		b := l.Authenticate(u, p)
		assert.True(t, b, "Failed authentication")
	}
}
