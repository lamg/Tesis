// Web interface for synchronizing UPR users databases
package tesis

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"html/template"
	"io/ioutil"
	h "net/http"
	"path"
)

const AuthHd = "Auth"

const (
	//Content files
	index  = "index.html"
	jquery = "jquery.js"
	info   = "info.html"
	//HTTPS server key files
	cert = "cert.pem"
	key  = "key.pem"
	//paths
	infoP = "/info"
	syncP = "/sync"
)

var (
	notFound = []byte("404 File not found")
	pkey     *rsa.PrivateKey
	auth     Authenticator
	db       DBManager
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
func ListenAndServe(u string, a Authenticator, d DBManager) {
	var bs []byte
	var e error
	auth, db = a, d
	// { auth,db:initialized }
	bs, e = ioutil.ReadFile(key)
	// { loaded.key.bs ≡ e = nil}
	if e == nil {
		// { loaded.key.bs }
		pkey, e = jwt.ParseRSAPrivateKeyFromPEM(bs)
		// { parsed.pkey ≡ e = nil }
		if e == nil {
			h.HandleFunc("/", indexH)
			h.HandleFunc("/s/", staticH)
			h.HandleFunc(infoP, infoH)
			h.HandleFunc(syncP, syncH)
			h.HandleFunc("/favicon.ico", h.NotFoundHandler().ServeHTTP)
			h.ListenAndServeTLS(u, cert, key, nil)
			// { serving tesis }
		}
		// { serving tesis ≡ e = nil }
	}
	return
}

// Handler of "/" path
func indexH(w h.ResponseWriter, r *h.Request) {
	if r.Method == h.MethodGet {
		serveHTML(w, index, nil)
	}
}

// Handler of "/s" path
func staticH(w h.ResponseWriter, r *h.Request) {
	var file string
	file = path.Base(r.URL.Path)
	file = path.Join("st", file)
	h.ServeFile(w, r, file)
}

// exists.p ≡ ⟨∃ i: i ∈ `ls`: i = p⟩
func serveHTML(w h.ResponseWriter, file string, d interface{}) {
	var e error
	var tm *template.Template
	var p string
	p = path.Join("st", file)
	tm, e = template.ParseFiles(p)
	// { e = nil ≡ exists.p ∧ parsed.tm }
	if e == nil {
		//TODO learn templates
		tm.Execute(w, d)
	} else {
		w.WriteHeader(404)
		w.Write([]byte(e.Error()))
	}
}

func infoH(w h.ResponseWriter, r *h.Request) {
	// { r.Method ∈ h.Method* }
	if r.Method == h.MethodPost {
		infoPost(w, r)
	} else if r.Method == h.MethodGet {
		infoGet(w, r)
	}
}

func infoPost(w h.ResponseWriter, r *h.Request) {
	//globals
	// pkey: *rsa.PrivateKey
	// auth: Authenticator
	//end
	var e error
	var user, pass string
	var v bool

	user, pass = r.FormValue("user"), r.FormValue("pass")
	v = auth.Authenticate(user, pass)
	if v {
		//user is authenticated
		var u *User
		var t *jwt.Token
		var js string
		u = &User{UserName: user}
		t = jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), u)
		js, e = t.SignedString(pkey)

		if e == nil {
			var ck *h.Cookie
			ck = &h.Cookie{Name: AuthHd, Value: js}
			h.SetCookie(w, ck)
			//cookie contains authentication token
			writeInfo(w, user)
			// { written.userInfo ≢ error }
		}
	} else {
		//TODO delete, handle cookie
		h.SetCookie(w, &h.Cookie{Name: AuthHd, Value: ""})
	}
	if e != nil {
		w.Write([]byte(e.Error()))
	}
	//r.Body.Close()
}

func infoGet(w h.ResponseWriter, r *h.Request) {
	//globals
	//p: *rsa.PublicKey
	//q: DBManager
	//end
	var t *jwt.Token
	var e error
	t, e = parseToken(r, &pkey.PublicKey)
	if e == nil && t.Valid {
		var clm jwt.MapClaims
		var us string

		// { t.Claims is a jwt.MapClaims }
		clm = t.Claims.(jwt.MapClaims)
		us = clm["user"].(string)
		writeInfo(w, us)
	} else {
		w.WriteHeader(401)
	}
}

func writeInfo(w h.ResponseWriter, user string) {
	var inf *Info
	var e error
	inf, e = db.UserInfo(user)
	// { loaded.inf ≡ e = nil }
	if e == nil {
		// { loaded.inf }
		serveHTML(w, info, inf)
	} else {
		// { ¬loaded.inf }
		w.Write([]byte(e.Error()))
		w.WriteHeader(500)
	}
	// { written.inf ≢ written.(e.Error()) ≡ e ≠ nil }
}

func parseToken(r *h.Request, p *rsa.PublicKey) (t *jwt.Token, e error) {
	var ck *h.Cookie
	ck, e = r.Cookie(AuthHd)
	t, e = jwt.Parse(ck.Value,
		func(x *jwt.Token) (a interface{}, d error) {
			a, d = p, nil
			return
		})
	return
}

func syncH(w h.ResponseWriter, r *h.Request) {
	// { r.Method ∈ h.Method* }
	if r.Method == h.MethodPost {
		// { sync Info parsed }
		// { sync Info processed }
	}
}
