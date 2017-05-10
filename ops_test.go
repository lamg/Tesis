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
		e = ss.SyncPend(rcp, u)
	}
	assert.NoError(t, e)
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
