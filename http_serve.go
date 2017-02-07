package tesis

import (
	"encoding/json"
	"github.com/gorilla/mux"
	//"io"
	//"net"
	h "net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc h.HandlerFunc
}

var routes = []Route{
	Route{"Root", "GET", "/", rootH},
	Route{"Root", "POST", "/", authH},
}

func serveTLS() (e error) {
	var (
		hp = ":10443"
		ce = "cert.pem"
		ke = "key.pem"
	)
	hr := mux.NewRouter()
	for _, i := range routes {
		hr.Methods(i.Method).
			Path(i.Pattern).
			Name(i.Name).
			Handler(i.HandlerFunc)
	}
	e = h.ListenAndServeTLS(hp, ce, ke, hr)
	return
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

type Credentials struct {
	user string
	pass string
}

func authH(w h.ResponseWriter, r *h.Request) {
	c, d := &Credentials{}, json.NewDecoder(r.Body)
	if e := d.Decode(&c); e == nil {
		var m []byte
		if r := auth(c.user, c.pass); r {
			m = []byte("OK")
		} else {
			m = []byte("¡Error!")
			w.WriteHeader(401)//401 is HTTP auth failed code
		}
		w.Write(m)
	}
}
