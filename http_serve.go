// Web portal for University of Pinar del Río
package tesis

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
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
	index  = "index"
	jquery = "jquery.js"
	//HTTPS server key files
	cert = "cert.pem"
	key  = "key.pem"
	//paths
	authP = "/a/auth"
	infoP = "/a/info"
)

var (
	notFound = []byte("404 File not found")
)

// Creates a new instance of HTTPPortal and starts serving.
//  u: URL to serve
//  a: User authentication interface
//  q: Database manager interface
// The directory where the program is executed must have
// the following structure:
//  (. cert.pem key.pem
//    (st index.html index.js
//        dash.html dash.js
//        jquery.js))
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
				Route{"Root", "GET", "/", pagesH},
				Route{"RootFiles", "GET", "/{file}", pagesH},
				Route{"Scripts", "GET", "/s/{file}", scriptH},
				Route{"Auth", "POST", authP,
					func(w h.ResponseWriter, r *h.Request) {
						authH(w, r, pk, a)
					},
				},
				Route{"Info", "GET", infoP,
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
				TLSNextProto:   map[string]func(*h.Server, *tls.Conn, h.Handler){}, //deactivate HTTP/2
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

// Handler of "/" path
func pagesH(w h.ResponseWriter, r *h.Request) {
	var file string
	file = mux.Vars(r)["file"]
	if file == "" {
		file = index
	}
	file = file + ".html"
	file = path.Join("st", file)

	// { exists file ~file~ in cwd }
	serveHTML(w, file)
}

// Handler of "/s" path
func scriptH(w h.ResponseWriter, r *h.Request) {
	//TODO doesn't load jquery.js
	var file string
	file = mux.Vars(r)["file"]
	file = path.Join("st", file)
	serveFile(w, file)
}

func serveFile(w h.ResponseWriter, file string) {
	var e error
	var bs []byte
	bs, e = ioutil.ReadFile(file)
	if e == nil {
		w.Write(bs)
	} else {
		w.WriteHeader(404)
		w.Write(notFound)
	}
}

func serveHTML(w h.ResponseWriter, file string) {
	var e error
	var tm *template.Template
	tm, e = template.ParseFiles(file)
	if e == nil {
		tm.Execute(w, nil)
	} else {
		w.WriteHeader(404)
		w.Write(notFound)
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
