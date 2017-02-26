// Web portal for University of Pinar del Río
package tesis

import (
	"context"
	"crypto/rsa"
	"encoding/json"

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
	//paths
	authP = "/auth"
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
				Route{"Root", "GET", "/",
					func(w h.ResponseWriter, r *h.Request) {
						rootH(w, r, &pk.PublicKey)
					},
				},
				Route{"Auth", "POST", authP,
					func(w h.ResponseWriter, r *h.Request) {
						authH(w, r, pk, a)
					},
				},
				Route{"Static content", "GET", "/st/{file}", stH},
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

func rootH(w h.ResponseWriter, r *h.Request, p *rsa.PublicKey) {
	var t *jwt.Token
	var e error
	t, e = parseToken(r, p)
	if e == nil && t.Valid {
		//user is already authenticated
		//dashboard page presented
	} else {
		//user is not authenticathed
		var t *template.Template
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
		//login page presented
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
	var c *Credentials
	var e error
	var m, bd []byte
	var ok = []byte("OK")
	var err = []byte("¡Error!")

	c = new(Credentials)
	bd, _ = ioutil.ReadAll(r.Body)
	//println(string(bd))
	if e = json.Unmarshal(bd, c); e == nil {
		var v bool
		v = a.Authenticate(c.User, c.Pass)
		if v {
			//user is authenticated
			var u *User
			var t *jwt.Token
			var js string

			u = &User{UserName: c.User}
			t = jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), u)
			js, e = t.SignedString(p)
			if e == nil {
				w.Header().Set(AuthHd, js)
				m = ok
				//the header contains authentication token
			} else {
				m = err
			}
		} else {
			//user authentication failed
			m = err
		}
	} else {
		m = []byte(e.Error())
	}
	w.Write(m)
	//r.Body.Close()
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
