package tesis

import (
	"crypto/tls"

	"fmt"

	"io/ioutil"
	h "net/http"
)

type PortalUser struct {
	client *h.Client
	url    string
	ck     *h.Cookie
}

func NewPortalUser(url string) (p *PortalUser) {
	cfg := &tls.Config{InsecureSkipVerify: true}
	tr := &h.Transport{TLSClientConfig: cfg}
	cl := &h.Client{Transport: tr}
	p = &PortalUser{client: cl, url: url}
	return
}

func (p *PortalUser) Auth(c *Credentials) (a bool, e error) {
	a = false
	var r *h.Response
	var u string
	u = fmt.Sprintf("https://%s%s", p.url, infoP)
	r, e = p.client.PostForm(u,
		map[string][]string{
			"user": []string{c.User},
			"pass": []string{c.Pass},
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
	var u string
	var q *h.Request
	if p.ck == nil {
		e = fmt.Errorf("Auth failed")
	}
	if e == nil {
		u = fmt.Sprintf("https://%s%s", p.url, infoP)
		q, e = h.NewRequest("GET", u, nil)
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

	r, e = p.client.Get(fmt.Sprintf("https://%s", p.url))
	if e == nil {
		b, e = ioutil.ReadAll(r.Body)
		if e == nil {
			s = string(b)
		}
	}
	return
}
