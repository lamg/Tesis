package tesis

import (
	"fmt"
)

func (s *DummyManager) Candidates() (a []Diff, e error) {
	a = s.cs
	//iterate DB and filter comparing with AD
	return
}

func (s *DummyManager) Synchronize(user string, a []Diff) (e error) {
	var r []int
	var ex bool
	r, ex = make([]int, 0, len(a)), true
	for i := 0; ex && i != len(a); i++ {
		var x, y int
		x, y, ex = 0, len(s.cs), false
		for x != y {
			ex = EqDBRecord(&a[i].DBRec, &s.cs[x].DBRec)
			if !ex {
				x = x + 1
			} else {
				x, r = y, append(r, x)
			}
		}
		if !ex {
			e = fmt.Errorf("Elemento %v no pertenece a candidatos a ser sincronizados", a[i])
		}
	}
	for i, _ := range r {
		s.cs = append(s.cs[:r[i]], s.cs[r[i]+1:]...)
	}
	return
}

type DummyAuth struct {
}

func (d *DummyAuth) Authenticate(u, p string) (b bool) {
	b = u == p
	return
}

type DummyManager struct {
	cs []Diff
}

func NewDummyManager() (m *DummyManager) {
	m = &DummyManager{
		cs: []Diff{
			Diff{
				LDAPRec:  DBRecord{Id: "0", IN: "8901191122", Name: "LUIS"},
				DBRec:    DBRecord{Id: "0", IN: "8901191122", Name: "Luis"},
				Exists:   true,
				Mismatch: true,
				Src:      "SIGENU",
			},
			Diff{
				DBRec:    DBRecord{Id: "1", IN: "9001191122", Name: "Coco"},
				Exists:   false,
				Mismatch: false,
				Src:      "SIGENU",
			},
			Diff{
				LDAPRec:  DBRecord{Id: "1", IN: "9001191122", Name: "Coco"},
				Exists:   true,
				Mismatch: false,
				Src:      "SIGENU",
			},
		},
	}
	return
}

func (m *DummyManager) UserInfo(u string) (inf *Info, e error) {
	var cs []Diff
	var re []Change
	cs, e = m.Candidates()
	if e == nil {
		re = make([]Change, 0)
		inf = &Info{Name: u, Matches: cs, Record: re}
	}
	//TODO populate inf with more meaningful information
	return
}
