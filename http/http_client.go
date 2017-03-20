package http

import (
	"crypto/tls"
	"fmt"

	"io/ioutil"
	h "net/http"
)

type PortalUser struct {
	client *h.Client
	ck     *h.Cookie
	auInf  string
	index  string
}

func NewPortalUser(url string) (p *PortalUser) {
	cfg := &tls.Config{InsecureSkipVerify: true}
	tr := &h.Transport{TLSClientConfig: cfg}
	cl := &h.Client{Transport: tr}
	ai := fmt.Sprintf("https://%s%s", url, dashP)
	in := fmt.Sprintf("https://%s", url)
	p = &PortalUser{client: cl, auInf: ai, index: in}
	return
}

func (p *PortalUser) Auth(user, pass string) (a bool, e error) {
	a = false
	var r *h.Response

	r, e = p.client.PostForm(p.auInf,
		map[string][]string{
			"user": []string{user},
			"pass": []string{pass},
		})
	if e == nil {
		if r.StatusCode == 200 {
			var cs []*h.Cookie
			cs = r.Cookies()
			if len(cs) == 1 {
				p.ck = cs[0]
				//token string stored
				a = true
			} else {
				e = fmt.Errorf("Cantidad de cookies (%d) incorrecta", len(cs))
			}
		} else {
			e = fmt.Errorf("%s", r.Status)
		}
	}
	return
}

func (p *PortalUser) Info() (s string, e error) {
	var q *h.Request
	if p.ck == nil {
		e = fmt.Errorf("Auth failed")
	}
	if e == nil {
		q, e = h.NewRequest("GET", p.auInf, nil)
	}
	if e == nil {
		q.AddCookie(p.ck)
		var rp *h.Response
		//create a request with appropiate header
		rp, e = p.client.Do(q)
		if e == nil && rp.StatusCode == 200 {
			var b []byte
			b, e = ioutil.ReadAll(rp.Body)
			if e == nil {
				s = string(b)
			}
		}
	}
	return
}

func (p *PortalUser) Index() (s string, e error) {
	var r *h.Response
	var b []byte

	r, e = p.client.Get(p.index)
	if e == nil {
		b, e = ioutil.ReadAll(r.Body)
		if e == nil {
			s = string(b)
		}
	}
	return
}
