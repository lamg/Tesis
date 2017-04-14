package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/lamg/tesis"
	a "github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	h "net/http/httptest"
	"testing"
)

var j string //json web token
var ke error

func TestA(t *testing.T) {
	var e error
	var bs []byte
	db = &tesis.DummyManager{}
	bs, ke = ioutil.ReadFile("key.pem")
	// { loaded.key.bs ≡ e = nil }
	if a.NoError(t, e) {
		// { loaded.key.bs }
		pkey, ke = jwt.ParseRSAPrivateKeyFromPEM(bs)
	}
	a.NoError(t, ke)
}

func TestAuth(t *testing.T) {
	var e error
	var r *h.ResponseRecorder
	var bs []byte
	if a.NoError(t, ke) {
		var cr *tesis.Credentials
		r = h.NewRecorder()
		cr = &tesis.Credentials{"a", "a"}
		bs, e = json.Marshal(cr)
	}
	var q *http.Request
	if a.NoError(t, e) {
		var rd io.Reader
		rd = bytes.NewReader(bs)
		q, e = http.NewRequest("POST", "", rd)
	}
	if a.NoError(t, e) {
		authH(r, q)
		j = r.Header().Get(AuthHd)
	}
	a.NotEmpty(t, j, "r:%v", r)
}

func TestUinf(t *testing.T) {
	var q *http.Request
	var e error
	if a.NoError(t, ke) {
		q, e = http.NewRequest("GET", "", nil)
	}
	var r *h.ResponseRecorder
	if a.NoError(t, e) {
		q.Header.Add(AuthHd, j)
		r = h.NewRecorder()
		uinfH(r, q)
	}
	var ui *tesis.UserInfo
	ui = new(tesis.UserInfo)
	e = ProcRes(r, ui)
	if a.NoError(t, e, "r=%s", string(r.Body.Bytes())) {
		t.Log(ui)
	}
}

func TestRecr(t *testing.T) {
	var e error
	var bs []byte
	if a.NoError(t, ke) {
		var pn *tesis.PageN
		pn = &tesis.PageN{PageN: 0}
		bs, e = json.Marshal(pn)
	}
	var q *http.Request
	if a.NoError(t, e) {
		var rd io.Reader
		rd = bytes.NewReader(bs)
		q, e = http.NewRequest("POST", "", rd)
	}
	var r *h.ResponseRecorder
	if a.NoError(t, e) {
		r = h.NewRecorder()
		q.Header.Add(AuthHd, j)
		recrH(r, q)
	}
	var pc *tesis.PageC
	pc = new(tesis.PageC)
	e = ProcRes(r, pc)
	if a.NoError(t, e) {
		t.Log(pc)
	}
}

func ProcRes(r *h.ResponseRecorder, v interface{}) (e error) {
	if r.Code == http.StatusOK {
		e = json.Unmarshal(r.Body.Bytes(), v)
	} else if r.Code == http.StatusBadRequest {
		var er *tesis.Error
		var eu error
		eu = json.Unmarshal(r.Body.Bytes(), er)
		if eu == nil {
			e = fmt.Errorf(er.Message)
		}
	} else {
		e = fmt.Errorf("Unknown code %d", r.Code)
	}
	return
}
