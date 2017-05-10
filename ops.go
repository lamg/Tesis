package tesis

import (
	"time"
)

func (ss *StateSys) SyncPend(r RecordReceptor,
	u string) (e error) {
	var i int
	var ps []Diff
	var chg Change
	if ss.UsrAct == nil {
		// { user u has no activity }
		ss.UsrAct = make(map[string]*Activity)
		ss.UsrAct[u] = &Activity{
			Record:   make([]Change, 0),
			Proposed: make([]Diff, 0),
		}
	}
	i, ps = 0, ss.UsrAct[u].Proposed
	chg = Change{
		SRec: make([]Diff, 0),
		FRec: make([]Diff, 0),
	}
	for i != len(ps) {
		if ps[i].Exists && ps[i].Mismatch {
			// { ps.i is an inconsistency }
			e = r.Update(ps[i].LDAPRec.Name, &ps[i].DBRec)
		} else if ps[i].Exists && !ps[i].Mismatch {
			// { ps.i is a deletion }
			e = r.Delete(ps[i].LDAPRec.Name)
		} else if !ps[i].Exists {
			// { ps.i is an addition }
			e = r.Create(ps[i].LDAPRec.Name, &ps[i].DBRec)
		}
		if e == nil {
			chg.SRec = append(chg.SRec, chgDiff(ps[i]))
		} else {
			chg.FRec = append(chg.FRec, chgDiff(ps[i]))
		}
		i = i + 1
	}
	chg.Time = time.Now()
	ss.UsrAct[u].Proposed = convEqDiff(
		delSuc(convDiffEq(ss.UsrAct[u].Proposed),
			convDiffEq(chg.SRec)))
	ss.UsrAct[u].Record = append(ss.UsrAct[u].Record, chg)
	return
}

func delSuc(pr, sc []Eq) (r []Eq) {
	r = make([]Eq, 0, len(pr)-len(sc))
	for _, j := range pr {
		var i int
		i = 0
		for i != len(sc) && !sc[i].Equals(j) {
			i = i + 1
		}
		if i == len(sc) {
			r = append(r, j)
		}
	}
	// { r = pr - sc }
	return
}

func chgDiff(d Diff) (r Diff) {
	r = Diff{
		LDAPRec: DBRecord{Id: d.LDAPRec.Id},
		DBRec:   DBRecord{Id: d.DBRec.Id},
	}
	return
}

func convDiffEq(d []Diff) (r []Eq) {
	r = make([]Eq, len(d))
	for i, j := range d {
		r[i] = j
	}
	return
}

func convEqDiff(q []Eq) (d []Diff) {
	d = make([]Diff, len(q))
	for i, j := range q {
		d[i] = j.(Diff)
	}
	return
}
