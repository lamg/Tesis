// Web portal for University of Pinar del Río
package tesis

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	_ "fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	h "net/http"
	"path"
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
	url    string
	routes []Route
	pkey   *rsa.PrivateKey
	srv    *h.Server
}

const (
	//Content files
	index  = "index.html"
	jquery = "st/jquery.js"
	//HTTPS server key files
	cert = "cert.pem"
	key  = "key.pem"
)

// Creates a new instance of HTTPPortal and starts serving.
//  u: URL to serve
//  a: User authentication interface
//  q: Database manager interface
// The directory where the program is executed must have
// the following structure:
//  (. index.html cert.pem key.pem (st jquery.js))
func NewHTTPPortal(u string, a Authenticator, q DBManager) (p *HTTPPortal, e error) {
	var bs []byte
	bs, e = ioutil.ReadFile(key)
	if e == nil {
		var pk *rsa.PrivateKey
		pk, e = jwt.ParseRSAPrivateKeyFromPEM(bs)
		if e == nil {
			var r []Route
			var hr *mux.Router
			r = []Route{
				Route{"Root", "GET", "/", rootH},
				Route{"Static content", "GET", "/st/{file}", stH},
				Route{"Auth", "POST", "/auth",
					func(w h.ResponseWriter, r *h.Request) {
						authH(w, r, pk, a)
					},
				},
				Route{"Info", "GET", "/info",
					func(w h.ResponseWriter, r *h.Request) {
						infoH(w, r, &pk.PublicKey, q)
					},
				},
			}
			hr = mux.NewRouter()
			for _, i := range r {
				hr.Methods(i.Method).
					Path(i.Pattern).
					Name(i.Name).
					Handler(i.HandlerFunc)
			}
			p = &HTTPPortal{url: u, routes: r, pkey: pk}
			p.srv = &h.Server{
				Addr:           u,
				Handler:        hr,
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}
		}
	}
	return
}

func (p *HTTPPortal) Serve() (e error) {
	e = p.srv.ListenAndServeTLS(cert, key)
	return
}

func (p *HTTPPortal) Shutdown(c context.Context) {
	p.srv.Shutdown(c)
}

func rootH(w h.ResponseWriter, r *h.Request) {
	var t *template.Template
	var e error
	// { exists file index in cwd }
	t, e = template.ParseFiles(index)
	if e == nil {
		t.Execute(w, &struct {
			Msg string
		}{
			Msg: "Hola",
		})
	} else {
		var m []byte
		m = []byte("Error loading index.html")
		w.Write(m)
	}
}

func stH(w h.ResponseWriter, r *h.Request) {
	var file string
	var bs []byte
	var e error
	file = mux.Vars(r)["file"]
	// { exists file ~file~ in cwd }
	bs, e = ioutil.ReadFile(path.Join("st", file))
	if e == nil {
		w.Write(bs)
	} else {
		w.WriteHeader(404)
		w.Write([]byte("404 File not found"))
	}
}

func authH(w h.ResponseWriter, r *h.Request, p *rsa.PrivateKey, a Authenticator) {
	var (
		c *Credentials
		d *json.Decoder
		e error
	)
	c, d = new(Credentials), json.NewDecoder(r.Body)
	if e = d.Decode(c); e == nil {
		var m []byte
		var r bool

		r = a.Authenticate(c.User, c.Pass)
		if r {
			//user is authenticated
			var (
				u  *User
				t  *jwt.Token
				js string
			)

			u = &User{UserName: c.User}
			t = jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), u)
			js, e = t.SignedString(p)
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

func infoH(w h.ResponseWriter, r *h.Request, p *rsa.PublicKey, q DBManager) {
	var t *jwt.Token
	var e error
	t, e = parseToken(r, p)
	if e == nil && t.Valid {
		var (
			inf *Info
			b   []byte
			clm jwt.MapClaims
			us  string
		)
		// { t.Claims is a jwt.MapClaims }
		clm = t.Claims.(jwt.MapClaims)
		us = clm["user"].(string)
		inf, e = q.UserInfo(us)
		if e == nil {
			b, e = json.Marshal(inf)
			if e == nil {
				w.Write(b)
			} else {
				w.Write([]byte("Marshal failed"))
				w.WriteHeader(500)
			}
		} else {
			w.Write([]byte("DB query failed"))
			w.WriteHeader(500)
		}
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
