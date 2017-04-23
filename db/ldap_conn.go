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

// New LDAP Authenticator connecting through TLS
//  lds: LDAP server address
//  sf: Suffix of user account (string after @)
//  ldp: LDAP server port
func NewLDAPAuth(lds, sf string, ldp int) (l *LDAPAuth, e error) {
	var adr string
	var cfg *tls.Config
	adr, cfg = fmt.Sprintf("%s:%d", lds, ldp),
		&tls.Config{InsecureSkipVerify: true}
	l = new(LDAPAuth)
	l.c, e = ldap.DialTLS("tcp", adr, cfg)
	l.sf = sf
	return
}

// Authenticate user
//  u: user (string before @ in user account)
//  p: password
func (l *LDAPAuth) Authenticate(u, p string) (b bool, e error) {
	e = l.c.Bind(u+l.sf, p)
	b = e == nil
	/*if b {
		var ms []string
		ea = Search(u, m.lauth.c)
		ms = ea.GetAttributeValues("memberOf")
		sort.Strings(ms)
		var r int
		var grp string
		grp = "OU=Gestion" //provisional
		//TODO definir como se va a marcar a los
		//usuarios y administradores de este programa
		//en el directorio activo
		r = sort.SearchStrings(ms, grp)
		b = r != len(ms) && ms[r] == grp
	}
	// { u belongs to synchronizers group or synchronizers
	// admin group}
	*/
	return
}

func (l *LDAPAuth) UserInfo(u string) (f *tesis.UserInfo,
	e error) {
	var n *ldap.Entry
	n, e = Search(u, l.c)
	if e == nil {
		f = &tesis.UserInfo{Name: n.GetAttributeValue("cn")}
	}
	return
}

func (l *LDAPAuth) Close() (e error) {
	l.c.Close()
	return
}

func Search(u string, c *ldap.Conn) (n *ldap.Entry, e error) {
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
	if e == nil && len(r.Entries) == 1 {
		n = r.Entries[0]
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
