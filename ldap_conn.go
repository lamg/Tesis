package tesis

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap"
)

// u: user
// p: password
func conn(u, p string) (c *ldap.Conn, e error) {
	var (
		//ldap server address
		lds = "ad.upr.edu.cu"
		//account suffix
		sf = "@upr.edu.cu"
		//ldap server port
		ldp = 636
	)
	
	adr, cfg, acc := fmt.Sprintf("%s:%d", lds, ldp),
		&tls.Config{InsecureSkipVerify: true},
		u+sf
	if c, e = ldap.DialTLS("tcp", adr, cfg); e == nil {
		e = c.Bind(acc, p)
	}
	return
}

// closes the connection
func search(u string, c *ldap.Conn) (n []*ldap.EntryAttribute, e error) {
	var (
		baseDN = "dc=upr,dc=edu,dc=cu"
		scope = ldap.ScopeWholeSubtree
		deref = ldap.NeverDerefAliases
		sizel = 0
		timel = 0
		tpeol = false //TypesOnly
		filter = fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", u)
		attrs = []string{}
		conts []ldap.Control = nil //[]Control
	)
	
	s := ldap.NewSearchRequest(baseDN, scope, deref,
		sizel, timel, tpeol, filter, attrs, conts)
	if r, e := c.Search(s); len(r.Entries) != 0 && e == nil {
		n = r.Entries[0].Attributes
	} else {
		e = fmt.Errorf("La búsqueda de %s falló", u)
	}
	return
}
