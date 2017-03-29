package db

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap"
)

type LDAPAuth struct {
	c  *ldap.Conn
	sf string //suffix of user account (string after @)
}

// Authenticate user
//  u: user (string before @ in user account)
//  p: password
func (l *LDAPAuth) Authenticate(u, p string) (b bool) {
	var e error
	e = l.c.Bind(u+l.sf, p)
	b = e == nil
	return
}

// New LDAP Authenticator connecting through TLS
//  lds: LDAP server address
//  sf: Suffix of user account (string after @)
//  ldp: LDAP server port
func NewLDAPAuth(lds, sf string, ldp int) (l *LDAPAuth, e error) {
	adr, cfg := fmt.Sprintf("%s:%d", lds, ldp),
		&tls.Config{InsecureSkipVerify: true}
	l = new(LDAPAuth)
	l.c, e = ldap.DialTLS("tcp", adr, cfg)
	l.sf = sf
	return
}

func (l *LDAPAuth) Close() (e error) {
	l.c.Close()
	return
}
