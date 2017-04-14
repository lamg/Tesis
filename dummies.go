package tesis

import (
	"fmt"
)

type DummyManager struct {
	cs []Diff
	pr []Diff
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
		pr: make([]Diff, 0),
	}
	return
}

func (m *DummyManager) UserInfo(u string) (inf *UserInfo, e error) {
	inf = &UserInfo{Name: u}
	return
}

func (d *DummyManager) Authenticate(u, p string) (b bool) {
	b = u == p
	return
}

func (d *DummyManager) Record(u string, p int) (c *PageC, e error) {
	c = new(PageC)
	return
}

func (d *DummyManager) Propose(u string, p []Diff) (e error) {
	d.pr = append(d.pr, p...)
	return
}

func (d *DummyManager) Pending(u string, p int) (c *PageD, e error) {
	c = &PageD{Total: 1, PageN: 1, DiffP: d.pr}
	return
}

func (d *DummyManager) Synchronize() (e error) {
	d.cs, e = RmEq(d.cs, d.pr)
	//save removed to record
	return
}

func RmEq(l, a []Diff) (p []Diff, e error) {
	var r []int
	var ex bool
	r, ex = make([]int, 0, len(a)), true
	for i := 0; ex && i != len(a); i++ {
		var x, y int
		x, y, ex = 0, len(l), false
		for x != y {
			ex = EqDBRecord(&a[i].DBRec, &l[x].DBRec)
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
	// { r contains the indexes in l of all elements in a
	// ≢ exists an element in a not in l }
	p = make([]Diff, 0, len(a))
	var c int
	c = 0
	for i := 0; e == nil && i != len(l); i++ {
		if c != len(r) && i == r[c] {
			c++
		} else {
			p = append(p, l[i])
		}
	}
	// { p contains elements in l not in a ≢
	// exists an element in a not in l }
	return
}