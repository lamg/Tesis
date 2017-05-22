package db

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/lamg/tesis"
	"strings"
)

const (
	IN                = "employeeID" //identity number
	userPrincipalName = "userPrincipalName"
	mail              = "mail"
	sAMAccountName    = "sAMAccountName"
	streetAddress     = "streetAddress"
	telephoneNumber   = "telephoneNumber"
	DN                = "distinguishedName"
	name              = "name"
	CN                = "cn"
	displayName       = "displayName"
)

func (l *LDAPAuth) Create(dn string,
	d *tesis.DBRecord) (e error) {
	var rq *ldap.AddRequest

	rq = ldap.NewAddRequest(dn)
	rq.Attribute(IN, []string{d.IN})
	rq.Attribute(displayName, []string{d.Name})
	rq.Attribute(streetAddress, []string{d.Addr})
	rq.Attribute(telephoneNumber, []string{d.Tel})
	e = l.c.Add(rq)
	return
}

func (l *LDAPAuth) Update(dn string,
	d *tesis.DBRecord) (e error) {
	// dn is us distinguished name
	var rq *ldap.ModifyRequest
	rq = ldap.NewModifyRequest(dn)
	rq.Replace(IN, []string{d.IN})
	rq.Replace(displayName, []string{d.Name})
	rq.Replace(streetAddress, []string{d.Addr})
	rq.Replace(telephoneNumber, []string{d.Tel})
	e = l.c.Modify(rq)
	return
}

func (l *LDAPAuth) Delete(dn string) (e error) {
	//var rq *ldap.DelRequest
	//change distinguished name to move the record
	//to another part of the LDAP tree
	//var rq *ldap.ModifyRequest
	//rq = ldap.NewModifyRequest(dn)
	//TODO
	// read backwards dn, is equal to ndn until
	// last DC, then OU=_Usuarios change to OU=_Graduados
	// or OU=_Baja according student status
	//rq.Add(DN, nDN)
	//e = l.c.Modify(rq)
	e = fmt.Errorf("Not implemented")
	return
}

type LDAPAuth struct {
	c     *ldap.Conn
	sf    string //suffix of user account (string after @)
	limit int
}

func NewLDAPProv(u, p, lda string, t int) (r tesis.RecordProvider, e error) {
	//set a limit
	var l *LDAPAuth
	const sf = "@upr.edu.cu" //account suffix
	l, e = NewLDAPAuth(lda, sf)
	var b bool
	if e == nil {
		b, e = l.Authenticate(u, p)
	}
	if e == nil && !b {
		e = fmt.Errorf("Falló al autenticar")
	}
	if e == nil {
		l.limit = t
		r = l
	}
	return
}

// New LDAP Authenticator connecting through TLS
//  lds: LDAP server address
//  sf: Suffix of user account (string after @)
//  ldp: LDAP server port
func NewLDAPAuth(lda, sf string) (l *LDAPAuth, e error) {
	var cfg *tls.Config
	cfg = &tls.Config{InsecureSkipVerify: true}
	l = new(LDAPAuth)
	l.c, e = ldap.DialTLS("tcp", lda, cfg)
	l.sf = sf
	return
}

// Authenticate user
//  u: user (string before @ in user account)
//  p: password
func (l *LDAPAuth) Authenticate(u, p string) (b bool, e error) {
	e = l.c.Bind(u+l.sf, p)
	b = e == nil
	return
}

func (l *LDAPAuth) Name() (s string) {
	s = "LDAP"
	return
}

func (l *LDAPAuth) UserInfo(u string) (f *tesis.UserInfo,
	e error) {
	var n *ldap.Entry
	var filter string
	var atts []string
	filter, atts =
		fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))",
			u),
		[]string{"cn"}
	n, e = SearchOne(filter, atts, l.c)
	if e == nil {
		f = &tesis.UserInfo{
			Name:     n.GetAttributeValue("cn"),
			UserName: u,
		}
	}
	return
}

func (l *LDAPAuth) UserRecord(us string) (d *tesis.DBRecord, e error) {
	var ats []string
	var flr string
	var n *ldap.Entry
	flr, ats = fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))", us),
		[]string{displayName, DN, IN, streetAddress, telephoneNumber}
	n, e = SearchOne(flr, ats, l.c)
	if e == nil {
		d = &tesis.DBRecord{
			Id:   n.GetAttributeValue(DN),
			IN:   n.GetAttributeValue(IN),
			Name: n.GetAttributeValue(displayName),
			Addr: n.GetAttributeValue(streetAddress),
			Tel:  n.GetAttributeValue(telephoneNumber),
		}
	}
	return
}

func (l *LDAPAuth) Close() (e error) {
	l.c.Close()
	return
}

func Search(u string, c *ldap.Conn) (av []string, e error) {
	var filter = fmt.Sprintf("(&(objectClass=user)(displayName=%s))",
		u)
	var attrs = []string{}
	var n *ldap.Entry
	n, e = SearchOne(filter, attrs, c)
	if e == nil {
		av = make([]string, 0, len(n.Attributes))
		for i := range n.Attributes {
			av = append(av, fmt.Sprintf("%s: %v",
				n.Attributes[i].Name,
				n.Attributes[i].Values))
		}
	}
	return
}

func (l *LDAPAuth) Records() (us []tesis.DBRecord, e error) {
	var f string
	var a []string
	f, a = "(&(objectCategory=person)(objectClass=user))",
		[]string{displayName, DN, IN, streetAddress, telephoneNumber}
	var n []*ldap.Entry
	n, e = SearchFilter(f, a, l.c)
	if e == nil && l.limit >= 0 && l.limit <= len(n) {
		n = n[:l.limit]
	}
	us = make([]tesis.DBRecord, 0, len(n))
	for _, i := range n {
		var r tesis.DBRecord
		r = tesis.DBRecord{
			Name: i.GetAttributeValue(displayName),
			Id:   strings.Join(i.GetAttributeValues(DN), ","),
			IN:   i.GetAttributeValue(IN),
			Addr: i.GetAttributeValue(streetAddress),
			Tel:  i.GetAttributeValue(telephoneNumber),
		}
		us = append(us, r)
	}
	return
}

func SearchDN(user string, c *ldap.Conn) (dn string, e error) {
	var n *ldap.Entry
	var filter string
	var atts []string
	filter, atts =
		fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))",
			user),
		[]string{DN}
	n, e = SearchOne(filter, atts, c)
	if e == nil {
		dn = strings.Join(n.GetAttributeValues(DN), ",")
	}
	return
}

func SearchOne(f string, ats []string, c *ldap.Conn) (n *ldap.Entry, e error) {
	var ns []*ldap.Entry
	ns, e = SearchFilter(f, ats, c)
	if e == nil {
		if len(ns) == 1 {
			n = ns[0]
		} else {
			e = fmt.Errorf("Result length = %d", len(ns))
		}
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
	} else if e == nil {
		n = r.Entries
	}
	return
}
