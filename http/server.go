// Web interface for synchronizing UPR users databases
package http

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lamg/tesis"
	"io/ioutil"
	"log"
	h "net/http"
)

const (
	//Content files
	AuthHd = "Auth"
	//paths
	authP = "/api/auth"
	recrP = "/api/recr"
	propP = "/api/prop"
	pendP = "/api/pend"
	chckP = "/api/chck"
	revpP = "/api/revp"
	rootP = "/"
	publP = "/public/"
)

var (
	notFound = []byte("Archivo no encontrado")
	notAuth  = []byte("No autenticado")
	pkey     *rsa.PrivateKey
	db       tesis.DBManager
)

// Creates a new instance of HTTPPortal and starts serving.
//  u: URL to serve
//  a: User authentication interface
//  d: Database manager interface
//  cr: cert.pem path
//  ky: key.pem path
//  h.DefaultServer: h.Server,!
func ListenAndServe(u string, d tesis.DBManager,
	cr, ky string) {
	var bs []byte
	var e error
	db = d
	// { auth,db,fs:initialized }
	bs, e = ioutil.ReadFile(ky)
	// { loaded.key.bs ≡ e = nil }
	if e == nil {
		// { loaded.key.bs }
		pkey, e = jwt.ParseRSAPrivateKeyFromPEM(bs)
	}
	// { parsed.pkey ≡ e = nil }
	if e == nil {
		h.HandleFunc(authP, authH)

		h.HandleFunc(recrP, recrH)
		h.HandleFunc(propP, propH)
		h.HandleFunc(pendP, pendH)
		h.HandleFunc(chckP, chckH)
		h.HandleFunc(revpP, revpH)
		// { API handlers set }

		h.Handle(rootP, h.FileServer(h.Dir("public")))
		// { static files handler set }

		h.ListenAndServeTLS(u, cr, ky, nil)
	}
	// { started.server ≡ e = nil }
	if e != nil {
		log.Print(e.Error())
	}
	return
}

// { db ≠ nil ∧ pkey ≠ nil }
func authH(w h.ResponseWriter, r *h.Request) {
	var e error
	var cr *tesis.Credentials
	var bs []byte

	if r.Method == h.MethodPost {
		bs, e = ioutil.ReadAll(r.Body)
		r.Body.Close()
	} else {
		e = errUnsMeth(r.Method, authP)
	}
	if e == nil {
		cr = new(tesis.Credentials)
		e = json.Unmarshal(bs, cr)
	}
	if e == nil {
		var v bool
		v, e = db.Authenticate(cr.User, cr.Pass)
		if !v {
			e = fmt.Errorf("Credenciales inválidas")
		}
		// { v ≡ registered.user ≡ e = nil }
	}
	var js string
	if e == nil {
		var u *tesis.User
		var t *jwt.Token
		u = &tesis.User{UserName: cr.User}
		t = jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), u)
		js, e = t.SignedString(pkey)
		// { signedString.js ≡ e = nil }
	}
	var ui *tesis.UserInfo
	if e == nil {
		ui, e = db.UserInfo(cr.User)
		ui.Token = js
	}
	// { infoLoaded ≡ e = nil }
	var rs []byte
	if e == nil {
		rs, e = json.Marshal(ui)
	}
	// { infoMarshaled ≡ e = nil}
	if e == nil {
		_, e = w.Write(rs)
	}
	writeError(w, e)
	// { writtenError ≢ writtenHeader }
}

func chckH(w h.ResponseWriter, r *h.Request) {
	var e error
	if r.Method != h.MethodGet {
		e = errUnsMeth(r.Method, chckP)
	}
	if e == nil {
		_, e = parseUserName(r, &pkey.PublicKey)
	}
	writeError(w, e)
}

func recrH(w h.ResponseWriter, r *h.Request) {
	var e error

	if r.Method != h.MethodPost {
		e = errUnsMeth(r.Method, recrP)
	}
	var us string
	if e == nil {
		us, e = parseUserName(r, &pkey.PublicKey)
	}
	var bs []byte
	if e == nil {
		bs, e = ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
	var pn *tesis.PageN
	if e == nil {
		pn = new(tesis.PageN)
		e = json.Unmarshal(bs, pn)
	}
	var pc *tesis.PageC
	if e == nil {
		pc, e = db.Record(us, pn.PageN)
	}
	var rs []byte
	if e == nil {
		rs, e = json.Marshal(pc)
	}
	if e == nil {
		_, e = w.Write(rs)
	}
	writeError(w, e)
}

func propH(w h.ResponseWriter, r *h.Request) {
	var us string
	var e error
	var patch, post bool
	patch, post = r.Method == h.MethodPatch,
		r.Method == h.MethodPost
	if patch || post {
		us, e = parseUserName(r, &pkey.PublicKey)

	} else {
		e = errUnsMeth(r.Method, propP)
	}
	var bs []byte
	if e == nil {
		bs, e = ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
	if e == nil && patch {
		var sel []string
		e = json.Unmarshal(bs, &sel)
		if e == nil {
			e = db.Propose(us, sel)
		}
	} else if e == nil && post {
		var pn *tesis.PageN
		pn = new(tesis.PageN)
		e = json.Unmarshal(bs, pn)
		var pd *tesis.PageD
		if e == nil {
			pd, e = db.Proposed(us, pn.PageN)
		}
		var rs []byte
		if e == nil {
			rs, e = json.Marshal(pd)
		}
		if e == nil {
			_, e = w.Write(rs)
		}
	}
	writeError(w, e)
}

func pendH(w h.ResponseWriter, r *h.Request) {
	var e error
	_, e = parseUserName(r, &pkey.PublicKey)
	var bs []byte
	if e == nil {
		if r.Method == h.MethodPost {
			bs, e = ioutil.ReadAll(r.Body)
			r.Body.Close()
		} else {
			e = errUnsMeth(r.Method, pendP)
		}
	}
	var pn *tesis.PageN
	if e == nil {
		pn = new(tesis.PageN)
		e = json.Unmarshal(bs, pn)
	}
	var pd *tesis.PageD
	if e == nil {
		pd, e = db.Pending(pn.PageN)
	}
	var rs []byte
	if e == nil {
		rs, e = json.Marshal(pd)
	}
	if e == nil {
		_, e = w.Write(rs)
	}
	writeError(w, e)
}

func revpH(w h.ResponseWriter, r *h.Request) {
	var us string
	var e error

	if r.Method == h.MethodPatch {
		us, e = parseUserName(r, &pkey.PublicKey)
	} else {
		e = errUnsMeth(r.Method, propP)
	}
	var bs []byte
	if e == nil {
		bs, e = ioutil.ReadAll(r.Body)
		r.Body.Close()
	}
	var sel []string
	if e == nil {
		e = json.Unmarshal(bs, &sel)
	}
	if e == nil {
		db.RevertProp(us, sel)
	}
	writeError(w, e)
}

func errUnsMeth(method, path string) (e error) {
	e = fmt.Errorf("Método %s no soportado por %s",
		method, path)
	return
}

func parseUserName(r *h.Request, p *rsa.PublicKey) (us string, e error) {
	var ts string
	var t *jwt.Token
	ts = r.Header.Get(AuthHd)
	// { readHeader.jwt ≡ e = nil }
	t, e = jwt.Parse(ts,
		func(x *jwt.Token) (a interface{}, d error) {
			a, d = p, nil
			return
		})
	// { parsedToken.t ≡ e = nil }
	if e == nil {
		if t.Valid {
			var clm jwt.MapClaims
			var ok bool
			clm, ok = t.Claims.(jwt.MapClaims)
			if ok {
				us, ok = clm["user"].(string)
			} else {
				e = fmt.Errorf("Error en aserción de tipo jwt.Claims")
			}
			if !ok && e == nil {
				e = fmt.Errorf("Error en aserción de tipo string")
			}
		} else {
			e = fmt.Errorf("Token de autenticación inválido")
		}
	} else {
		e = fmt.Errorf("Token de autenticación malformado")
	}
	return
}

func writeError(w h.ResponseWriter, e error) {
	if e != nil {
		var in *tesis.Error
		var bs []byte
		in = &tesis.Error{Message: e.Error()}
		bs, e = json.Marshal(in)
		if e != nil {
			// precondition of json.Marshal is false
			// i.e. program is incorrect
			log.Panicf("Incorrect program: %s", e.Error())
		} else {
			w.WriteHeader(h.StatusBadRequest)
			_, e = w.Write(bs)
		}
	}
	// { writtenError ≡ e ≠ nil }
	if e != nil {
		log.Print(e.Error())
	}
}

/*
{ A ≢ ¬R }
if R → { A } S { B }
  ¬R → skip
fi
{ B ≢ ¬R }
*/
