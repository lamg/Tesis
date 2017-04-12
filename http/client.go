package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/lamg/tesis"
	"io"
	"io/ioutil"
	"log"
	h "net/http"
)

type PortalUser struct {
	client *h.Client
	jwt    string
	auInf  string
	index  string
}

func NewPortalUser(url string) (p *PortalUser) {
	cfg := &tls.Config{InsecureSkipVerify: true}
	tr := &h.Transport{TLSClientConfig: cfg}
	cl := &h.Client{Transport: tr}

	in := fmt.Sprintf("https://%s", url)
	p = &PortalUser{client: cl, index: in}
	return
}

func (p *PortalUser) Auth(user, pass string) (a bool, e error) {
	a = false
	var r *h.Response
	var cr *tesis.Credentials
	var bs []byte
	cr = &tesis.Credentials{User: user, Pass: pass}
	bs, e = json.Marshal(cr)
	if e == nil {
		var rd io.Reader
		rd = bytes.NewReader(bs)
		r, e = p.client.Post(p.index+authP, "application/json", rd)
	} else {
		log.Panicf("Incorrect program: %s", e.Error())
	}
	if e == nil {
		if r.StatusCode == 200 {
			//get Auth header
			p.jwt = r.Header.Get(AuthHd)
			//token string stored
			if p.jwt != "" {
				a = true
			} else {
				bs, e = ioutil.ReadAll(r.Body)
				if e == nil {
					e = fmt.Errorf("Cabecera Auth vacía, %s", string(bs))
				}
			}
		} else {
			e = fmt.Errorf("Status = %s", r.Status)
		}
	}
	return
}

func (p *PortalUser) Sync() (s string, e error) {
	//marshall Info to JSON and send to server
	//to use syncH handler
	var r *h.Response
	var b []byte
	var acs []tesis.Diff
	var rd io.Reader
	acs = []tesis.Diff{
		tesis.Diff{
			LDAPRec:  tesis.DBRecord{Id: "0", IN: "8901191122", Name: "LUIS"},
			DBRec:    tesis.DBRecord{Id: "0", IN: "8901191122", Name: "Luis"},
			Exists:   true,
			Mismatch: true,
			Src:      "SIGENU",
		},
	}
	b, e = json.Marshal(acs)
	rd = bytes.NewReader(b)
	var q *h.Request
	if p.jwt == "" {
		e = fmt.Errorf("Auth failed")
	}
	if e == nil {
		q, e = h.NewRequest("POST", p.index+syncP, rd)
	}
	if e == nil {
		q.Header.Set(AuthHd, p.jwt)
		q.Header.Set("content-type", "application/json")
		r, e = p.client.Do(q)
	}
	var bs []byte
	if e == nil {
		bs, e = ioutil.ReadAll(r.Body)
	}
	if e == nil {
		s = string(bs)
		e = r.Body.Close()
	}
	// { syncRespStr.s ≡ e = nil }
	return
}
