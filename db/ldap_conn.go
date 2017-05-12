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
	CN                = "cn"
)

type LDAPRecp struct {
	c *ldap.Conn
}

func NewLDAPRecp(adr, u,
	p string) (r tesis.RecordReceptor, e error) {
	var cfg *tls.Config
	var l *LDAPRecp
	cfg, l = &tls.Config{InsecureSkipVerify: true},
		new(LDAPRecp)

	l.c, e = ldap.DialTLS("tcp", adr, cfg)
	r = l
	return
}

func (l *LDAPRecp) Create(dn string,
	d *tesis.DBRecord) (e error) {
	var rq *ldap.AddRequest

	rq = ldap.NewAddRequest(dn)
	rq.Attribute(IN, []string{d.IN})
	rq.Attribute(CN, []string{d.Name})
	rq.Attribute(streetAddress, []string{d.Addr})
	rq.Attribute(telephoneNumber, []string{d.Tel})
	e = l.c.Add(rq)
	return
}

func (l *LDAPRecp) Update(dn string,
	d *tesis.DBRecord) (e error) {
	var rq *ldap.ModifyRequest
	rq = ldap.NewModifyRequest(dn)
	rq.Add(IN, []string{d.IN})
	rq.Add(CN, []string{d.Name})
	rq.Add(streetAddress, []string{d.Addr})
	rq.Add(telephoneNumber, []string{d.Tel})
	e = l.c.Modify(rq)
	return
}

func (l *LDAPRecp) Delete(dn string) (e error) {
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

func (l *LDAPAuth) Name() (s string) {
	s = "LDAP"
	return
}

func (l *LDAPAuth) UserInfo(u string) (f *tesis.UserInfo,
	e error) {
	var n []*ldap.Entry
	var filter string
	var atts []string
	filter, atts =
		fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s))",
			u),
		[]string{"cn"}
	n, e = SearchFilter(filter, atts, l.c)
	if len(n) == 0 {
		e = fmt.Errorf("Busqueda de información fallo en AD")
	}
	if e == nil {
		f = &tesis.UserInfo{
			Name:     n[0].GetAttributeValue("cn"),
			UserName: u,
		}
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
		filter                = fmt.Sprintf("(&(objectClass=user)(cn=%s))", u)
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

func (l *LDAPAuth) Records() (us []tesis.DBRecord, e error) {
	var f string
	var a []string
	f, a = "(&(objectCategory=person)(objectClass=user))",
		[]string{"cn", DN, IN}
	var n []*ldap.Entry
	n, e = SearchFilter(f, a, l.c)
	if e == nil && l.limit >= 0 && l.limit <= len(n) {
		n = n[:l.limit]
	}
	us = make([]tesis.DBRecord, 0, len(n))
	for _, i := range n {
		var r tesis.DBRecord
		r = tesis.DBRecord{}
		var ln int
		var b bool
		ln = len(i.Attributes)
		b = ln >= 1
		if b {
			r.Name = i.Attributes[0].Values[0]
			b = ln >= 2
		}
		if b {
			r.Id = strings.Join(i.Attributes[1].Values, ",")
			b = ln >= 3
		}
		if b {
			r.IN = i.Attributes[2].Values[0]
		}
		if r.Name != "" && r.Id != "" {
			us = append(us, r)
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
