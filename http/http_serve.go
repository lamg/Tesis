// Web interface for synchronizing UPR users databases
package http

import (
	"crypto/rsa"
	"crypto/tls"
	"github.com/dgrijalva/jwt-go"
	"github.com/lamg/tesis"
	"html/template"
	"io/ioutil"
	"log"
	h "net/http"
	"path"
)

const (
	AuthHd = "auth"
	//Content files
	jquery = "jquery.js"
	//HTTPS server key files
	cert = "cert.pem"
	key  = "key.pem"
	//paths
	dashP = "/dash"
	syncP = "/sync"
)

var (
	notFound = []byte("Archivo no encontrado")
	notAuth  = []byte("No autenticado")
	tms      *template.Template
	fTms     = []string{"st/index.html", "st/dash.html"}
	pkey     *rsa.PrivateKey
	auth     tesis.Authenticator
	db       tesis.DBManager
	indexTm  *template.Template
	dashTm   *template.Template
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
//        jquery.js util.js))
func ListenAndServe(u string, a tesis.Authenticator, d tesis.DBManager) {
	var bs []byte
	var e error
	h.DefaultClient.Transport = &h.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) h.RoundTripper),
	}
	// { disabled.HTTP2 }
	auth, db = a, d
	// { auth,db:initialized }
	bs, e = ioutil.ReadFile(key)
	// { loaded.key.bs ≡ e = nil }
	if e == nil {
		// { loaded.key.bs }
		pkey, e = jwt.ParseRSAPrivateKeyFromPEM(bs)
	}
	// { parsed.pkey ≡ e = nil }
	if e == nil {
		// exists.p ≡ ∃.`ls`.(=p)
		tms, e = template.ParseFiles(fTms...)
		// { ∀.fTms.exists ∧ parsed.tm ≡ e = nil }
	}
	if e == nil {
		h.HandleFunc("/", indexH)
		h.HandleFunc("/s/", staticH)
		h.HandleFunc(dashP, dashH)
		h.HandleFunc(syncP, syncH)
		h.HandleFunc("/favicon.ico", h.NotFoundHandler().ServeHTTP)
		h.ListenAndServeTLS(u, cert, key, nil)
	}
	// { started.server ≡ e = nil }
	if e != nil {
		log.Print(e.Error())
	}
	return
}

// Handler of "/" path
func indexH(w h.ResponseWriter, r *h.Request) {
	if r.Method == h.MethodGet {
		tms.ExecuteTemplate(w, path.Base(fTms[0]), nil)
	}
}

// Handler of "/s" path
func staticH(w h.ResponseWriter, r *h.Request) {
	var file string
	file = path.Base(r.URL.Path)
	file = path.Join("st", file)
	h.ServeFile(w, r, file)
}

func dashH(w h.ResponseWriter, r *h.Request) {
	// { r.Method ∈ h.Method* }
	if r.Method == h.MethodPost {
		dashPost(w, r)
		// { written.UserName ∧ written.Cookie ≢ writtenError }
	} else if r.Method == h.MethodGet {
		dashGet(w, r)
		//
	}
}

func dashPost(w h.ResponseWriter, r *h.Request) {
	//globals
	// pkey: *rsa.PrivateKey,?
	// auth: Authenticator,?
	// AuthHd: string,?
	//end
	var e error
	var user, pass string
	var v bool
	e, v = r.ParseForm(), false
	if e == nil {
		user, pass = r.PostFormValue("user"), r.PostFormValue("pass")
		v = auth.Authenticate(user, pass)
	}
	// { v ≡ registered.user }
	if v {
		var u *tesis.User
		var t *jwt.Token
		var js string
		u = &tesis.User{UserName: user}
		t = jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), u)
		js, e = t.SignedString(pkey)
		// { signedString.js ≡ e = nil }
		if e == nil {
			var ck *h.Cookie
			ck = &h.Cookie{Name: AuthHd, Value: js}
			h.SetCookie(w, ck)
			writeInfo(w, u.UserName)
			// { written.(u.UserName) ∧ written.ck }
		}
	} else {
		w.Write(notAuth)
	}
	if e != nil {
		w.Write([]byte(e.Error()))
	}
	// { written.(u.UserName) ∧ written.ck ≢ written.(e.Error()) }
}

func dashGet(w h.ResponseWriter, r *h.Request) {
	//globals
	//p: *rsa.PublicKey,?
	//end
	var t *jwt.Token
	var e error
	t, e = parseToken(r, &pkey.PublicKey)
	// { e = nil ∧ t.Valid ≡ auth.(user.t) }
	if e == nil && t.Valid {
		var clm jwt.MapClaims
		var us string

		// { t.Claims: jwt.MapClaims }
		clm = t.Claims.(jwt.MapClaims)
		us = clm["user"].(string)
		writeInfo(w, us)
		// { writtenInfo.us }
	} else {
		// { e ≠ nil ∨ ¬t.Valid }
		w.Write(notAuth)
		w.WriteHeader(401)
	}
	// { (writtenInfo.us ≢ written.notAuth }
}

func writeInfo(w h.ResponseWriter, user string) {
	// globals
	// db: DBManager,?
	// fTms: []string,?
	// tms: *template.Template,?
	// end
	var inf *tesis.Info
	var e error
	inf, e = db.UserInfo(user)
	// { loaded.inf ≡ e = nil }
	if e == nil {
		// { loaded.inf }
		tms.ExecuteTemplate(w, path.Base(fTms[1]), inf)
		// { written.inf }
	} else {
		// { ¬loaded.inf }
		w.Write([]byte(e.Error()))
		w.WriteHeader(500)
		// { written.(e.Error()) }
	}
	// { written.inf ≢ writtenError }
}

func parseToken(r *h.Request, p *rsa.PublicKey) (t *jwt.Token, e error) {
	var ck *h.Cookie
	ck, e = r.Cookie(AuthHd)
	// { readCookie.ck ≡ e = nil }
	if e == nil {
		t, e = jwt.Parse(ck.Value,
			func(x *jwt.Token) (a interface{}, d error) {
				a, d = p, nil
				return
			})
	}
	// { parsedToken.t ≡ e = nil }
	return
}

func syncH(w h.ResponseWriter, r *h.Request) {
	// { r.Method ∈ h.Method* }
	if r.Method == h.MethodPost {
		// { sync Info parsed }
		// { sync Info processed }
	}
}
