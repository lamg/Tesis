// Web portal for University of Pinar del Río
package tesis

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	_ "fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	h "net/http"
	"time"
)

const AuthHd = "Auth"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc h.HandlerFunc
}

type HTTPPortal struct {
	url, cert, key string
	routes         []Route
	pkey           *rsa.PrivateKey
}

// Creates a new instance of HTTPPortal and starts serving
//  u: URL to serve
//  c: cert.pem path
//  k: key.pem path
//  a: User authentication interface
func NewHTTPPortal(u, c, k string, a Authenticator) (p *HTTPPortal, e error) {
	var pk *rsa.PrivateKey
	pk, e = rsa.GenerateKey(rand.Reader, 2048)
	r := []Route{
		Route{"Root", "GET", "/", rootH},
		Route{"Auth", "POST", "/",
			func(w h.ResponseWriter, r *h.Request) {
				authH(w, r, pk, a)
			},
		},
		Route{"Info", "GET", "/info",
			func(w h.ResponseWriter, r *h.Request) {
				infoH(w, r, &pk.PublicKey)
			},
		},
	}

	if e == nil {
		p = &HTTPPortal{url: u, cert: c, key: k, routes: r,
			pkey: pk}
		hr := mux.NewRouter()
		for _, i := range r {
			hr.Methods(i.Method).
				Path(i.Pattern).
				Name(i.Name).
				Handler(i.HandlerFunc)
		}
		go func() {
			e = h.ListenAndServeTLS(u, c, k, hr)
		}()
		time.Sleep(1000000 * time.Nanosecond)
	}
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

func authH(w h.ResponseWriter, r *h.Request, p *rsa.PrivateKey, a Authenticator) {
	c, d := &Credentials{}, json.NewDecoder(r.Body)
	if e := d.Decode(&c); e == nil {
		var m []byte
		if r := a.Authenticate(c.User, c.Pass); r {
			//user is authenticated
			t := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), c)
			js, e := t.SignedString(p)
			if e == nil {
				w.Header().Set(AuthHd, js)
				m = []byte("OK")
				//the header contains authentication token
			} else {
				m = []byte("Error")
				w.WriteHeader(401)
			}
		} else {
			m = []byte("¡Error!")
			w.WriteHeader(401) //401 is HTTP auth failed code
		}
		w.Write(m)
	}
}

func infoH(w h.ResponseWriter, r *h.Request, p *rsa.PublicKey) {
	t, e := parseToken(r, p)
	if e == nil && t.Valid {
		//get info for t.Claims["user"]
		//write info to the response
	} else {
		w.WriteHeader(401)
	}
}

func parseToken(r *h.Request, p *rsa.PublicKey) (t *jwt.Token, e error) {
	var js string
	js = r.Header.Get(AuthHd)
	t, e = jwt.Parse(js,
		func(x *jwt.Token) (a interface{}, d error) {
			a, d = p, nil
			return
		})
	return
}
