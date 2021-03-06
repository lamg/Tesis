package tesis

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSyncPend(t *testing.T) {
	var bs []byte
	bs = []byte(ssJSON)
	var e error
	var ss *StateSys
	ss = new(StateSys)
	e = json.Unmarshal(bs, ss)
	var rcp RecordReceptor
	var u string
	u, rcp = "luis.mendez", NewDRCP(t)
	if assert.NoError(t, e) {
		ss.UsrAct = make(map[string]*Activity)
		ss.UsrAct[u] = &Activity{Proposed: ss.Pending}
		ss.Pending = make([]Diff, 0)
		var pr Reporter
		pr = NewTRpr(t)
		e = ss.SyncPend(rcp, u, pr)
	}
	assert.NoError(t, e)
	if assert.True(t, ss.UsrAct != nil &&
		ss.UsrAct[u] != nil) {
		var dr DBRecord
		dr = ss.UsrAct[u].Record[0].SRec[0].DBRec
		if !assert.True(t, dr.Equals(dbRec)) {
			t.Log(dr)
		}
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

var dbRec = DBRecord{
	Id:   "91742be:1501970c670:-3d",
	IN:   "95120923357",
	Name: "Claudia Crúz Labrador",
	Addr: "Km 10 Carretera Viñales, CPA Isidro Barre do, Viñalesdo",
	Tel:  "",
}
