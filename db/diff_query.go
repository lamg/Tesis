package db

import (
	"github.com/lamg/tesis"
)

func PDiff(dbProv, ldProv tesis.RecordProvider,
	rp tesis.Reporter) (ds []tesis.Diff, e error) {
	var st, us []tesis.DBRecord
	st, e = dbProv.Records()
	if e == nil {
		us, e = ldProv.Records()
	}
	var g, h, x, y []tesis.Sim
	if e == nil {
		x, y = convSim(st), convSim(us)
		_, g, h, _ = tesis.DiffSym(x, y, rp)
		// { ¬ (g,h contain equal couples) }
	}
	var k, l []tesis.DBRecord
	if e == nil {
		k, l = convDBR(g), convDBR(h)
		/*ds = make([]tesis.Diff, 0, len(j)+len(k)+len(m))
		for _, jx := range j {
			ds = append(ds, tesis.Diff{
				DBRec:    jx,
				Src:      dbProv.Name(),
				Exists:   false,
				Mismatch: false,
			})
		}*/
		// { ds contains LDAP additions }
		for ix, jx := range k {
			ds = append(ds, tesis.Diff{
				DBRec:    jx,
				LDAPRec:  l[ix],
				Src:      dbProv.Name(),
				Exists:   true,
				Mismatch: true,
			})
		}
		// { ds contains LDAP mismatches }
		/*for _, jx := range m {
			ds = append(ds, tesis.Diff{
				LDAPRec:  jx,
				Src:      dbProv.Name(),
				Exists:   true,
				Mismatch: false,
			})
		}*/
		// { ds contains LDAP deletions }
		// { ds contains all pending operations for dbProv }
	}
	return
}

func convSim(s []tesis.DBRecord) (r []tesis.Sim) {
	r = make([]tesis.Sim, len(s))
	for i, j := range s {
		r[i] = j
	}
	return
}

func convDBR(s []tesis.Sim) (r []tesis.DBRecord) {
	r = make([]tesis.DBRecord, len(s))
	for i, j := range s {
		r[i] = j.(tesis.DBRecord)
	}
	return
}
