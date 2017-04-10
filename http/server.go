// Web interface for synchronizing UPR users databases
package http

import (
	"crypto/rsa"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lamg/tesis"
	"io"
	"io/ioutil"
	"log"
	h "net/http"
	"path"
)

const (
	//Content files
	AuthHd = "Auth"
	//paths
	syncP = "/api/sync"
	authP = "/api/auth"
)

var (
	notFound = []byte("Archivo no encontrado")
	notAuth  = []byte("No autenticado")
	pkey     *rsa.PrivateKey
	auth     tesis.Authenticator
	db       tesis.DBManager
)

type ServFS struct {
	Cert, Key string
}

// Creates a new instance of HTTPPortal and starts serving.
//  u: URL to serve
//  a: User authentication interface
//  q: Database manager interface
//  f: Files needed to run the server
//  db: []AccMatch,!
//  h.DefaultServer: h.Server,!
func ListenAndServe(u string, a tesis.Authenticator, d tesis.DBManager, f *ServFS) {
	var bs []byte
	var e error
	h.DefaultClient.Transport = &h.Transport{
		TLSNextProto: make(map[string]func(authority string, c *tls.Conn) h.RoundTripper),
	}
	// { disabled.HTTP2 }
	auth, db = a, d
	// { auth,db,fs:initialized }
	bs, e = ioutil.ReadFile(f.Key)
	// { loaded.key.bs ≡ e = nil }
	if e == nil {
		// { loaded.key.bs }
		pkey, e = jwt.ParseRSAPrivateKeyFromPEM(bs)
	}
	// { parsed.pkey ≡ e = nil }
	if e == nil {
		h.HandleFunc(syncP, syncH)
		h.HandleFunc(authP, authH)
		h.HandleFunc("/", indexH)
		h.ListenAndServeTLS(u, f.Cert, f.Key, nil)
	}
	// { started.server ≡ e = nil }
	if e != nil {
		log.Print(e.Error())
	}
	return
}

// Handler of "/" path
func indexH(w h.ResponseWriter, r *h.Request) {
	var e error
	if r.Method == h.MethodGet {
		var file, ext string
		file = path.Base(r.URL.Path)
		if file == "/" {
			file = "index"
		}
		ext = path.Ext(file)
		if ext == ".js" {
			w.Header().Set("content-type", "application/javascript")
		} else if ext == ".css" {
			w.Header().Set("content-type", "text/css")
		} else if ext == "" {
			file = file + ".html"
			w.Header().Set("content-type", "text/html")
		}
		h.ServeFile(w, r, file)

	} else {
		e = fmt.Errorf("Método %s no soportado por /", r.Method)
	}
	writeError(w, e)
}

func writeError(w h.ResponseWriter, e error) {
	if e != nil {
		var in *tesis.Info
		var bs []byte
		in = &tesis.Error{Message: e.Error()}
		bs, e = json.Marshal(in)
		if e != nil {
			// precondition of json.Marshal is false
			// i.e. program is incorrect
			log.Panicf("Incorrect program: %s", e.Error())
		} else {
			w.WriteHeader(400)
			_, e = w.Write(bs)
		}
	}
	// { writtenError ≡ e ≠ nil }
	if e != nil {
		log.Print(e.Error())
	}
}

func authH(w h.ResponseWriter, r *h.Request) {
	var e error
	var cr *tesis.Credentials
	var bs []byte
	if r.Method == h.MethodPost {
		bs, e = ioutil.ReadAll(r.Body)
	} else {
		e = fmt.Errorf("Method %s not supported by /a/auth", r.Method)
	}
	if e == nil {
		cr = new(tesis.Credentials)
		e = json.Unmarshal(bs, cr)
	}
	var js string
	if e == nil {
		var v bool
		v = auth.Authenticate(cr.User, cr.Pass)
		// { v ≡ registered.user }
		if v {
			var u *tesis.User
			var t *jwt.Token
			u = &tesis.User{UserName: cr.User}
			t = jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), u)
			js, e = t.SignedString(pkey)
			// { signedString.js ≡ e = nil }
		} else {
			e = fmt.Errorf("Credenciales inválidas")
		}
	}
	if e == nil {
		w.Header().Set(AuthHd, js)
		// { header set }
	}
	writeError(w, e)
	// { writtenError ≢ writtenCookie }
}

func parseUserName(r *h.Request, p *rsa.PublicKey) (us string, e error) {
	var ck *h.Cookie
	var t *jwt.Token
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
	if e == nil {
		if t.Valid {
			var clm jwt.MapClaims
			// TODO what can go wrong here?
			// { t.Claims: jwt.MapClaims }
			clm = t.Claims.(jwt.MapClaims)
			us = clm["user"].(string)
		} else {
			e = fmt.Errorf("Token no válido")
		}
	}
	return
}

func syncH(w h.ResponseWriter, r *h.Request) {
	var e error
	var us string
	us, e = parseUserName(r, &pkey.PublicKey)
	if e == nil {
		if r.Method == h.MethodPost {
			e = syncPost(w, r.Body, us)
		} else if r.Method == h.MethodGet {
			e = syncGet(w, us)
		} else {
			e = fmt.Errorf("%s no soportado en /sync", r.Method)
		}
	}
	writeError(w, e)
}

func syncPost(w h.ResponseWriter, rc io.Reader, us string) (e error) {
	var bs []byte
	// { r.Method = h.MethodPost ∧ validJWT ≡ e = nil }
	bs, e = ioutil.ReadAll(rc)
	// { read.bs ≡ e = nil }
	var acs []tesis.AccMatch
	if e == nil {
		// { read.bs }
		e = json.Unmarshal(bs, &acs)
	}
	// { jsonRep.bs.acs ≡ e = nil }
	if e == nil {
		e = db.Synchronize(us, acs)
	}
	var cs []tesis.AccMatch
	if e == nil {
		cs, e = db.Candidates()
	}
	var rs []byte
	if e == nil {
		rs, e = json.Marshal(cs)
	}
	if e == nil {
		_, e = w.Write(rs)
	}
	// { synchronized.acs ≡ e = nil }
	return
}

func syncGet(w h.ResponseWriter, user string) (e error) {
	// globals
	// db: DBManager,?
	// end
	var inf *tesis.Info
	var bs []byte
	inf, e = db.UserInfo(user)
	// { loaded.inf ≡ e = nil }
	if e == nil {
		// { loaded.inf }
		bs, e = json.Marshal(inf)
		// { written.inf }
	}
	if e == nil {
		w.Header().Set("content-type", "application/json")
		_, e = w.Write(bs)
		//sent relevant information as JSON, according user
	}
	// { written.inf ≡ e = nil }
	return
}

/*
{ A ≢ ¬R }
if R → { A } S { B }
  ¬R → skip
fi
{ B ≢ ¬R }
*/
