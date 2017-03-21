// Web interface for synchronizing UPR users databases
package http

import (
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lamg/tesis"
	"html/template"
	"io/ioutil"
	"log"
	h "net/http"
	"path"
)

const (
	//Content files
	AuthHd = "auth"
	//paths
	dashP = "/dash"
	syncP = "/sync"
)

var (
	notFound = []byte("Archivo no encontrado")
	notAuth  = []byte("No autenticado")
	tms      *template.Template
	pkey     *rsa.PrivateKey
	auth     tesis.Authenticator
	db       tesis.DBManager
	indexTm  *template.Template
	dashTm   *template.Template
	fs       *ServFS
)

type ServFS struct {
	JsFiles, Cert, Key string
	FTms               []string
}

// Creates a new instance of HTTPPortal and starts serving.
//  u: URL to serve
//  a: User authentication interface
//  q: Database manager interface
//  f: Files needed to run the server
func ListenAndServe(u string, a tesis.Authenticator, d tesis.DBManager, f *ServFS) {
	var bs []byte
	var e error
	h.DefaultClient.Transport = &h.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) h.RoundTripper),
	}
	// { disabled.HTTP2 }
	auth, db, fs = a, d, f
	// { auth,db,fs:initialized }
	bs, e = ioutil.ReadFile(fs.Key)
	// { loaded.key.bs ≡ e = nil }
	if e == nil {
		// { loaded.key.bs }
		pkey, e = jwt.ParseRSAPrivateKeyFromPEM(bs)
	}
	// { parsed.pkey ≡ e = nil }
	if e == nil {
		// exists.p ≡ ∃.`ls`.(=p)
		tms, e = template.ParseFiles(fs.FTms...)
		// { ∀.fTms.exists ∧ parsed.tm ≡ e = nil }
	}
	if e == nil {
		h.HandleFunc("/", indexH)
		h.HandleFunc("/s/", staticH)
		h.HandleFunc(dashP, dashH)
		h.HandleFunc(syncP, syncH)
		h.HandleFunc("/favicon.ico", h.NotFoundHandler().ServeHTTP)
		h.ListenAndServeTLS(u, fs.Cert, fs.Key, nil)
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
		tms.ExecuteTemplate(w, path.Base(fs.FTms[0]), nil)
	}
}

// Handler of "/s" path
func staticH(w h.ResponseWriter, r *h.Request) {
	var file string
	file = path.Base(r.URL.Path)
	file = path.Join(fs.JsFiles, file)
	h.ServeFile(w, r, file)
}

func dashH(w h.ResponseWriter, r *h.Request) {
	// { r.Method ∈ h.Method* }
	if r.Method == h.MethodPost {
		dashPost(w, r)
		// { written.UserName ∧ written.Cookie ≢ writtenError }
	} else if r.Method == h.MethodGet {
		dashGet(w, r)
		// {}
	} else {
		log.Printf("Method %s not supported\n", r.Method)
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
		log.Println(e.Error())
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
		tms.ExecuteTemplate(w, path.Base(fs.FTms[1]), inf)
		// { written.inf }
	} else {
		// { ¬loaded.inf }
		log.Println(e.Error())
	}
	// { written.inf ≢ loggedError }
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
	if e != nil {
		log.Println(e.Error())
	}
	return
}

func syncH(w h.ResponseWriter, r *h.Request) {
	var e error
	var bs []byte

	// { r.Method ∈ h.Method* }
	if r.Method == h.MethodPost {
		var t *jwt.Token
		t, e = parseToken(r, &pkey.PublicKey)
		if e == nil && !t.Valid {
			e = fmt.Errorf("Token no válido")
		}
	} else {
		e = fmt.Errorf("Método %s no soportado", r.Method)
	}
	if e == nil {
		bs, e = ioutil.ReadAll(r.Body)
	}
	// { read.bs ≡ e = nil }
	if e == nil {
		// { read.bs }
		var acs []tesis.AccMatch
		e = json.Unmarshal(bs, acs)
	}
	// { jsonRep.bs.acs ≡ e = nil }
	if e != nil {
		log.Println(e.Error())
	}
}
