package db

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/lamg/tesis"
)

const IN = "in" //identity number

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

func Search(u string, c *ldap.Conn) (n []*ldap.EntryAttribute, e error) {
	var (
		baseDN                = "dc=upr,dc=edu,dc=cu"
		scope                 = ldap.ScopeWholeSubtree
		deref                 = ldap.NeverDerefAliases
		sizel                 = 0
		timel                 = 0
		tpeol                 = false //TypesOnly
		filter                = fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", u)
		attrs                 = []string{}
		conts  []ldap.Control = nil //[]Control
		s      *ldap.SearchRequest
		r      *ldap.SearchResult
	)

	s = ldap.NewSearchRequest(baseDN, scope, deref,
		sizel, timel, tpeol, filter, attrs, conts)
	r, e = c.Search(s)
	if e == nil && len(r.Entries) != 0 {
		n = r.Entries[0].Attributes
	} else if e == nil {
		e = fmt.Errorf("La búsqueda de %s falló", u)
	}
	// { attributes.u.n ≡ e = nil }
	return
}

func (l *LDAPAuth) GetUsers() (us []tesis.DBRecord, e error) {
	var f string
	var a []string
	f, a = "(&(objectClass=user))", []string{"cn",
		"userPrincipalName", IN}
	var n []*ldap.Entry
	n, e = SearchFilter(f, a, l.c)
	us = make([]tesis.DBRecord, len(n))
	for _, i := range n {
		var r tesis.DBRecord
		r = tesis.DBRecord{
			Name: i.GetAttributeValue("cn"),
			Id:   i.GetAttributeValue("userPrincipalName"),
			IN:   i.GetAttributeValue(IN),
		}
		us = append(us, r)
	}
	return
}

func SearchFilter(f string, ats []string, c *ldap.Conn) (n []*ldap.Entry, e error) {
	var (
		baseDN                = "dc=upr,dc=edu,dc=cu"
		scope                 = ldap.ScopeWholeSubtree
		deref                 = ldap.NeverDerefAliases
		sizel                 = 0
		timel                 = 0
		tpeol                 = false //TypesOnly
		conts  []ldap.Control = nil   //[]Control
		s      *ldap.SearchRequest
		r      *ldap.SearchResult
	)
	s = ldap.NewSearchRequest(baseDN, scope, deref,
		sizel, timel, tpeol, f, ats, conts)
	r, e = c.Search(s)
	if e == nil && len(r.Entries) == 0 {
		e = fmt.Errorf("La búsqueda de %s falló", f)
	} else {
		n = r.Entries
	}
	return
}
