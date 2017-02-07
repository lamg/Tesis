package tesis

import (
	"encoding/json"
	"github.com/gorilla/mux"
	//"io"
	//"net"
	h "net/http"
	"time"
	"bytes"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc h.HandlerFunc
}

func rootH(w h.ResponseWriter, r *h.Request) {
	var (
		ct = "Content-Type"
		tp = "text/plain"
		cs = "charset"
		ut = "utf-8"
		ms = []byte("¡Hola Mundo!")
	)

	w.Header().Set(ct, tp)
	w.Header().Set(cs, ut)
	w.Write(ms)
}

func authH(w h.ResponseWriter, r *h.Request) {
	c, d := &Credentials{}, json.NewDecoder(r.Body)
	if e := d.Decode(&c); e == nil {
		var m []byte
		if r := auth(c.user, c.pass); r {
			m = []byte("OK")
		} else {
			m = []byte("¡Error!")
			w.WriteHeader(401) //401 is HTTP auth failed code
		}
		w.Write(m)
	}
}

type HTTPPortal struct {
	url, cert, key string
	authS          bool
	auth, index    Route
	client *h.Client
}

func NewHTTPPortal(u, c, k string, iR, aR Route) (p *HTTPPortal, e error) {
	cl := &h.Client{}
	p = &HTTPPortal{url: u, cert: c, key: k, index: iR,
		auth: aR, authS: false, client: cl}
	hr := mux.NewRouter()
	for _, i := range []Route {iR, aR} {
		hr.Methods(i.Method).
			Path(i.Pattern).
			Name(i.Name).
			Handler(i.HandlerFunc)
	}
	go func() {
		e = h.ListenAndServeTLS(u, c, k, hr)
	}()
	time.Sleep(1000000 * time.Nanosecond)
	return
}

func (p *HTTPPortal) Auth(c *Credentials) (x bool) {
	x = false
	b, e := json.Marshal(c)
	if e == nil {
		br := bytes.NewReader(b)
		r, e := p.client.Post(p.url, "application/json", br)
		if e == nil {
			x = r.StatusCode == 200
		}
	}
	return
}

func (p *HTTPPortal) Conversate() (c []Conversation) {
	
}
