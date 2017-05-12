package http

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lamg/tesis"
	dbs "github.com/lamg/tesis/db"
	a "github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	h "net/http/httptest"

	"testing"
	"time"
)

var j string //json web token
var local = "http://localhost"

func TestA(t *testing.T) {
	var e error
	var bs []byte
	db = &tesis.DummyManager{}
	bs, e = ioutil.ReadFile("key.pem")
	// { loaded.key.bs ≡ e = nil }
	if a.NoError(t, e) {
		// { loaded.key.bs }
		pkey, e = jwt.ParseRSAPrivateKeyFromPEM(bs)
	}
	a.NoError(t, e)
}

func TestAuth(t *testing.T) {
	var e error
	var r *h.ResponseRecorder
	var bs []byte
	e = errPkey()
	if a.NoError(t, e) {
		var cr *tesis.Credentials
		r = h.NewRecorder()
		cr = &tesis.Credentials{"a", "a"}
		bs, e = json.Marshal(cr)
	}
	var q *http.Request
	if a.NoError(t, e) {
		var rd io.Reader
		rd = bytes.NewReader(bs)
		q = h.NewRequest(http.MethodPost, local, rd)
		authH(r, q)
		j = r.Body.String()
	}
	var b, c bool
	b, c = j == "", r != nil && (r.Code == http.StatusBadRequest)
	a.False(t, b || c)
}

func TestAuth0(t *testing.T) {
	var e error
	var bs []byte
	e = errPkey()
	if a.NoError(t, e) {
		var cr *tesis.Credentials
		cr = &tesis.Credentials{"a", "b"}
		bs, e = json.Marshal(cr)
	}
	var rd io.Reader
	var er *tesis.Error
	if a.NoError(t, e) {
		var q *http.Request
		var r *h.ResponseRecorder
		rd, er, r = bytes.NewReader(bs), new(tesis.Error),
			h.NewRecorder()
		q = h.NewRequest(http.MethodPost, local, rd)
		authH(r, q)
		e = json.Unmarshal(r.Body.Bytes(), er)
	}
	a.NoError(t, e)
}

func TestAuth1(t *testing.T) {
	var e error
	e = errPkey()
	if a.NoError(t, e) {
		a.HTTPError(t, authH, http.MethodGet, local, nil)
	}
}

func TestUinf(t *testing.T) {
	var e error
	var q *http.Request
	e = errAuth()
	if a.NoError(t, e) {
		q, e = http.NewRequest(http.MethodGet, "", nil)
	}
	var r *h.ResponseRecorder
	if a.NoError(t, e) {
		q.Header.Add(AuthHd, j)
		r = h.NewRecorder()
		uinfH(r, q)
	}
	var ui *tesis.UserInfo
	if a.NoError(t, e) {
		ui = new(tesis.UserInfo)
		e = procRes(r, ui)
	}
	a.NoError(t, e)
}

func TestUinf0(t *testing.T) {
	var e error
	e = errAuth()
	if a.NoError(t, e) {
		a.HTTPError(t, uinfH, http.MethodConnect, local, nil)
	}
}

func TestRecr(t *testing.T) {
	var e error
	var bs []byte
	e = errAuth()
	if a.NoError(t, e) {
		var pn *tesis.PageN
		pn = &tesis.PageN{PageN: 0}
		bs, e = json.Marshal(pn)
	} else {
		e = fmt.Errorf("Auth failed %s", t.Name())
	}
	var q *http.Request
	if a.NoError(t, e) {
		var rd io.Reader
		rd = bytes.NewReader(bs)
		q, e = http.NewRequest("POST", "", rd)
	}
	var r *h.ResponseRecorder
	var pc *tesis.PageC
	if a.NoError(t, e) {
		r = h.NewRecorder()
		q.Header.Add(AuthHd, j)
		recrH(r, q)
		pc = new(tesis.PageC)
		e = procRes(r, pc)
	}
	a.NoError(t, e)
}

func TestRecr0(t *testing.T) {
	var e error
	e = errAuth()
	if a.NoError(t, e) {
		a.HTTPError(t, recrH, http.MethodConnect, local, nil)
	}
}

func TestProp(t *testing.T) {
	var e error
	var bs []byte
	e = errAuth()
	if a.NoError(t, e) {
		var ds []tesis.Diff
		ds = []tesis.Diff{
			tesis.Diff{
				LDAPRec:  tesis.DBRecord{Id: "0", IN: "8901191122", Name: "LUIS"},
				DBRec:    tesis.DBRecord{Id: "0", IN: "8901191122", Name: "Luis"},
				Exists:   true,
				Mismatch: true,
				Src:      "SIGENU",
			},
		}
		bs, e = json.Marshal(ds)
	} else {
		e = fmt.Errorf("Auth failed %s", t.Name())
	}
	var q *http.Request
	if a.NoError(t, e) {
		var rd io.Reader
		rd = bytes.NewReader(bs)
		q, e = http.NewRequest(http.MethodPatch, "", rd)
	}
	var r *h.ResponseRecorder
	if a.NoError(t, e) {
		r = h.NewRecorder()
		q.Header.Add(AuthHd, j)
		propH(r, q)
	}
	a.True(t, e == nil &&
		r.Code == http.StatusOK &&
		r.Body.Len() == 0)
}

func TestProp0(t *testing.T) {
	var e error
	e = errAuth()
	if a.NoError(t, e) {
		a.HTTPError(t, propH, http.MethodConnect, local, nil)
	}
}

func TestPend(t *testing.T) {
	var e error
	var bs []byte
	e = errAuth()
	if a.NoError(t, e) {
		var pn *tesis.PageN
		pn = &tesis.PageN{PageN: 0}
		bs, e = json.Marshal(pn)
	}
	var q *http.Request
	if a.NoError(t, e) {
		var rd io.Reader
		rd = bytes.NewReader(bs)
		q, e = http.NewRequest(http.MethodPost, "", rd)
	}
	var r *h.ResponseRecorder
	var pd *tesis.PageD
	if a.NoError(t, e) {
		r, pd = h.NewRecorder(), new(tesis.PageD)
		pendH(r, q)
		e = procRes(r, pd)
	}
	a.NoError(t, e)
}

func TestPend0(t *testing.T) {
	var e error
	e = errAuth()
	if a.NoError(t, e) {
		a.HTTPError(t, pendH, http.MethodConnect, local, nil)
	}
}

func TestFileServ(t *testing.T) {
	var q *http.Request
	var e error
	q, e = http.NewRequest(http.MethodGet, rootP, nil)
	var r *h.ResponseRecorder
	if a.NoError(t, e) {
		r = h.NewRecorder()
		http.FileServer(http.Dir(".")).ServeHTTP(r, q)
		a.True(t, r.Code == http.StatusOK, "Code: %d", r.Code)
	}
}

func procRes(r *h.ResponseRecorder, v interface{}) (e error) {
	if r.Code == http.StatusOK {
		e = json.Unmarshal(r.Body.Bytes(), v)
	} else if r.Code == http.StatusBadRequest {
		var er *tesis.Error
		er = new(tesis.Error)
		e = json.Unmarshal(r.Body.Bytes(), er)
		if e == nil {
			e = fmt.Errorf(er.Message)
		}
	} else {
		e = fmt.Errorf("Unknown code %d", r.Code)
	}
	// { v is the JSON value when http.StatusOK
	// ≡ e = nil ≢ e has msg sent when http.StatusBadRequest
	// ∨ e is unknown code error}
	return
}

func errPkey() (e error) {
	if pkey == nil {
		e = fmt.Errorf("Nil pkey")
	}
	return
}

func errAuth() (e error) {
	e = errPkey()
	if e == nil && j == "" {
		e = fmt.Errorf("Failed auth")
	}
	return
}

func TestServ(t *testing.T) {
	var d tesis.DBManager
	var y tesis.UserDB
	var jsonC, clAddr, srAddr string
	var e error
	y, jsonC, clAddr, srAddr = tesis.NewDummyManager(),
		"application/json", "https://localhost:10443",
		":10443"
	var rd io.Reader
	var wr io.Writer
	rd, wr = bytes.NewBufferString(ssJSON),
		bytes.NewBufferString("")
	d, e = dbs.NewUPRManager(rd, wr, y)
	if a.NoError(t, e) {
		go ListenAndServe(srAddr, d, "cert.pem", "key.pem")
	}
	time.Sleep(200 * time.Millisecond)
	// wait for the server to load
	//Test auth
	var tc *tls.Config
	tc = &tls.Config{InsecureSkipVerify: true}
	http.DefaultTransport = &http.Transport{
		TLSNextProto: make(map[string]func(authority string,
			c *tls.Conn) http.RoundTripper),
		// { disabledHTTP2 }
		TLSClientConfig: tc,
		// { disabledSelfSignedCertCheck }
	}
	// { configuredClient }

	var bs []byte
	var cr *tesis.Credentials
	cr = &tesis.Credentials{"a", "a"}
	bs, e = json.Marshal(cr)
	var rs *http.Response
	if a.NoError(t, e) {
		var rd io.Reader
		rd = bytes.NewReader(bs)
		rs, e = http.Post(clAddr+"/api/auth", jsonC, rd)
	}
	if a.NoError(t, e) {
		a.EqualValues(t, 200, rs.StatusCode)
	}
}

var ssJSON = `
{
   "pending": [
     {
			"ldapRec": {
				"id": "CN=Claudia Crúz Labrador,OU=4to,OU=MarxismoHistoria,OU=CRD,OU=Pregrado,OU=Estudiantes,OU=FEM,OU=Facultades,OU=_Usuarios,DC=upr,DC=edu,DC=cu",
				"in": "",
				"name": "Claudia Crúz Labrador",
				"addr": "",
				"tel": ""
			},
			"dbRec": {
				"id": "91742be:1501970c670:-3d",
				"in": "95120923357",
				"name": "Claudia Crúz Labrador",
				"addr": "Km 10 Carretera Viñales, CPA Isidro Barre do, Viñalesdo",
				"tel": ""
			},
			"src": "sigenu",
			"exists": true,
			"mismatch": true
		}
	],
	"usrAct": null
}

`
