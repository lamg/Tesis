package tesis

import (
	"fmt"
)

func (s *DummyManager) Candidates() (a []AccMatch, e error) {
	a = s.cs
	//iterate DB and filter comparing with AD
	return
}

func (s *DummyManager) Synchronize(user string, a []AccMatch) (e error) {
	var r []int
	var ex bool
	r, ex = make([]int, 0, len(a)), true
	for i := 0; ex && i != len(a); i++ {
		var x, y int
		x, y, ex = 0, len(s.cs), false
		for x != y {
			ex = EqAccMatch(&a[i], &s.cs[x])
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
	cs []AccMatch
}

func NewDummyManager() (m *DummyManager) {
	m = &DummyManager{
		cs: []AccMatch{
			AccMatch{DBId: "0", SrcIN: "8901191122"},
			AccMatch{DBId: "1", ADId: "3", ADName: "LUIS", SrcName: "Luis", SrcDB: "ASET"},
			AccMatch{DBId: "2", SrcIN: "9001091221"},
		},
	}
	return
}

func (m *DummyManager) UserInfo(u string) (inf *Info, e error) {
	var cs []AccMatch
	var re []Change
	cs, e = m.Candidates()
	if e == nil {
		re = make([]Change, 0)
		inf = &Info{Name: u, Matches: cs, Record: re}
	}
	//TODO populate inf with more meaningful information
	return
}
