package db

import (
	"github.com/lamg/tesis"
	"log"
)

func Sync(dbProv, ldProv tesis.RecordProvider) (ds []tesis.Diff, e error) {
	var st, us []tesis.DBRecord
	st, e = dbProv.Records()
	if e == nil {
		us, e = ldProv.Records()
	}
	var f, g, h, i, x, y []tesis.Sim
	if e == nil {
		st, us = st[:1000], us[:1000] //Test purposes
		x, y = convSim(st), convSim(us)
		f, g, h, i = tesis.DiffSym(x, y)
		// { ¬ (g,h contain equal couples) }
		log.Printf("f=%d,g=%d,h=%d,i=%d", len(f), len(g), len(h),
			len(i))
	}
	//TODO
	var j, k, l, m []tesis.DBRecord
	if e == nil {
		j, k, l, m = convDBR(f), convDBR(g), convDBR(h),
			convDBR(i)
		ds = make([]tesis.Diff, 0, len(j)+len(k)+len(m))
		log.Printf("ds=%d", len(ds))
		for _, jx := range j {
			ds = append(ds, tesis.Diff{
				DBRec:    jx,
				Src:      dbProv.Name(),
				Exists:   false,
				Mismatch: false,
			})
		}
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
		for _, jx := range m {
			ds = append(ds, tesis.Diff{
				LDAPRec:  jx,
				Src:      dbProv.Name(),
				Exists:   true,
				Mismatch: false,
			})
		}
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
