package tesis

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	h "net/http"
)

type PortalUser struct {
	client  *h.Client
	url, tk string
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
	var b []byte
	b, e = json.Marshal(c)
	if e == nil {
		var br io.Reader
		var r *h.Response
		var u string
		br = bytes.NewReader(b)
		u = fmt.Sprintf("https://%s/auth", p.url)
		r, e = p.client.Post(u, "application/json", br)
		if e == nil {
			if r.StatusCode == 200 {
				p.tk = r.Header.Get(AuthHd)
				var bd []byte
				bd, e = ioutil.ReadAll(r.Body)
				println(string(bd))
				//token string stored
				a = true
			} else {
				e = fmt.Errorf("%s", r.Status)
			}
		} else {

		}
	}
	return
}

func (p *PortalUser) Info() (inf *Info, e error) {
	var u string
	var q *h.Request

	u = fmt.Sprintf("https://%s/%s", p.url, "info")
	q, e = h.NewRequest("GET", u, nil)
	if e == nil {
		var rp *h.Response
		q.Header = map[string][]string{
			AuthHd: {p.tk},
		}
		//create a request with appropiate header
		rp, e = p.client.Do(q)
		if e == nil && rp.StatusCode == 200 {
			var d *json.Decoder
			inf = new(Info)
			d = json.NewDecoder(rp.Body)
			e = d.Decode(inf)
			rp.Body.Close()
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
