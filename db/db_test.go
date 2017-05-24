package db

import (
	"bytes"
	"github.com/lamg/tesis"
	a "github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestPageSlice(t *testing.T) {
	var s, r []interface{}
	s = make([]interface{}, 20)
	for i, _ := range s {
		s[i] = i
	}
	var n, l, m, ps int
	n = 3 // page size
	r, l, m, ps = pageSlice(s, n, 6)
	a.True(t, len(r) == 2 && l == 18 && m == 20 && ps == 7)
}

func TestPending(t *testing.T) {
	var um *UPRManager
	var e error
	var dm tesis.UserDB
	dm = &tesis.DummyManager{}
	var rwc io.ReadWriteCloser
	rwc = tesis.NewRWC(bytes.NewBufferString(ssJSON),
		bytes.NewBufferString(""))
	um, e = NewUPRManager(rwc, dm)
	var pd *tesis.PageD
	if a.NoError(t, e) {
		pd, e = um.Pending(0)
	}
	a.True(t, e == nil && pd != nil)
}

func TestPropose(t *testing.T) {
	var pr []string
	pr = []string{"91742be:1501970c670:-3d"}
	//each proposed exists in pending
	var um *UPRManager
	var e error
	var dm tesis.UserDB
	var user string
	dm, user = &tesis.DummyManager{}, "lamg"
	//FIXME no esta leyendo el archivo
	var rwc io.ReadWriteCloser
	rwc = tesis.NewRWC(bytes.NewBufferString(ssJSON),
		bytes.NewBufferString(""))
	um, e = NewUPRManager(rwc, dm)
	if a.NoError(t, e) {
		e = um.Propose(user, pr)
	}
	var pd *tesis.PageD
	if a.NoError(t, e) {
		pd, e = um.Pending(0)
	}
	var pp *tesis.PageD
	if a.NoError(t, e) && a.True(t, len(pd.DiffP) == 0,
		"len %d ≠ 0", len(pd.DiffP)) {
		pp, e = um.Proposed(user, 0)
		a.True(t, e == nil && pp.DiffP != nil &&
			len(pp.DiffP) == 1 &&
			pp.DiffP[0].DBRec.Id == pr[0])
	}

	//each proposed doesn't exists in pending
}

func TestRevertProp(t *testing.T) {
	var um *UPRManager
	var e error
	var dm tesis.UserDB
	var pr []string
	var user string
	dm = &tesis.DummyManager{}

	var rwc io.ReadWriteCloser
	rwc = tesis.NewRWC(bytes.NewBufferString(ssJSON),
		bytes.NewBufferString(""))
	um, e = NewUPRManager(rwc, dm)
	if a.NoError(t, e) {
		pr, user = []string{"91742be:1501970c670:-3d"}, "lamg"
	}
	e = um.Propose(user, pr)
	if a.NoError(t, e) {
		e = um.RevertProp(user, pr)
	}
	var pp *tesis.PageD
	if a.NoError(t, e) {
		pp, e = um.Proposed(user, 0)
	}
	var pd *tesis.PageD
	if a.NoError(t, e) {
		pd, e = um.Pending(0)
	}
	a.True(t, e == nil && pp.DiffP != nil &&
		len(pp.DiffP) == 0 && pd.DiffP != nil &&
		len(pd.DiffP) == 1 && pd.DiffP[0].DBRec.Id == pr[0])
}

func TestSymDiff0(t *testing.T) {
	var x, y []tesis.Diff
	x, y =
		[]tesis.Diff{
			tesis.Diff{
				LDAPRec:  tesis.DBRecord{Id: "0", IN: "8901191122", Name: "LUIS"},
				DBRec:    tesis.DBRecord{Id: "0", IN: "8901191122", Name: "Luis"},
				Exists:   true,
				Mismatch: true,
				Src:      "SIGENU",
			},
		},
		[]tesis.Diff{
			tesis.Diff{
				LDAPRec:  tesis.DBRecord{Id: "0", IN: "8901191122", Name: "LUIS"},
				DBRec:    tesis.DBRecord{Id: "0", IN: "8901191122", Name: "Luis"},
				Exists:   true,
				Mismatch: true,
				Src:      "SIGENU",
			},
		}
	a.True(t, x[0].DBRec.Equals(y[0].DBRec))
	a.True(t, x[0].Equals(y[0]))
	var u, v []tesis.Eq
	u, v = make([]tesis.Eq, len(x)), make([]tesis.Eq, len(y))
	u[0], v[0] = x[0], y[0]
	var w, z []tesis.Eq
	w, z = tesis.DiffInt(u, v)
	a.True(t, len(w) == 0 && len(z) == 1 &&
		z[0].Equals(x[0]))
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
