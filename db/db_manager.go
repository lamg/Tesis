package db

import (
	"encoding/json"
	"fmt"
	"github.com/lamg/tesis"
	"io"
	"io/ioutil"
)

type UPRManager struct {
	usrDB  tesis.UserDB
	steSys *tesis.StateSys
	// { JSONRepr.usrDt = contents.dtFile
	//   ≡ exists.dtFile }
	writer io.WriteCloser
	pgLen  int // { pgLen ≥ 0 }
}

func NewUPRManager(f io.ReadWriteCloser, a tesis.UserDB) (m *UPRManager, e error) {
	var bs []byte
	bs, e = ioutil.ReadAll(f)
	var ss *tesis.StateSys
	if e == nil {
		ss = new(tesis.StateSys)
		e = json.Unmarshal(bs, ss)
	}
	// { JSONRepr.ud = contents.dtFile
	//   ≡ exists.dtFile
	//   ≡ e = nil }
	if e == nil {
		m = &UPRManager{
			usrDB:  a,
			steSys: ss,
			writer: f,
			pgLen:  10,
		}
	}
	return
}

func (m *UPRManager) Authenticate(u, p string) (b bool, e error) {
	b, e = m.usrDB.Authenticate(u, p)
	if e == nil && b && m.steSys.UsrAct == nil {
		m.steSys.UsrAct = make(map[string]*tesis.Activity)
		m.steSys.UsrAct[u] = &tesis.Activity{
			Proposed: make([]tesis.Diff, 0),
			Record:   make([]tesis.Change, 0),
		}
	}
	return
}

func (m *UPRManager) UserInfo(u string) (n *tesis.UserInfo, e error) {
	n, e = m.usrDB.UserInfo(u)
	return
}

func (m *UPRManager) Record(u string, p int) (c *tesis.PageC, e error) {
	// { Authenticated.u }
	var r []tesis.Change
	r = m.steSys.UsrAct[u].Record

	var t []tesis.Change
	var a, b []interface{}
	a = make([]interface{}, len(r))
	for i, j := range r {
		a[i] = j
	}
	var ps int
	b, _, _, ps = pageSlice(a, m.pgLen, p)
	t = make([]tesis.Change, len(b))
	for i, j := range b {
		t[i] = j.(tesis.Change)
	}
	c = &tesis.PageC{Total: ps, PageN: p, ChangeP: t}
	return
}

func (m *UPRManager) Propose(u string, ds []string) (e error) {
	// { Authenticated.u }
	var d []tesis.Diff
	d = tesis.CreateDiff(ds)

	var f, g, h, l []tesis.Eq
	f, g = tesis.ConvDiffEq(m.steSys.Pending),
		tesis.ConvDiffEq(d)
	h, l = tesis.DiffInt(f, g)

	var k, n []tesis.Diff
	k, n = tesis.ConvEqDiff(h), tesis.ConvEqDiff(l)
	//{ k = 'k - d  ∧  n = 'k ∩ d }
	m.steSys.Pending = k
	m.steSys.UsrAct[u].Proposed = append(m.steSys.UsrAct[u].Proposed, n...)
	return
}

func (m *UPRManager) Close() (e error) {
	var bs []byte
	bs, e = json.MarshalIndent(m.steSys, "", "\t")
	if e == nil {
		_, e = m.writer.Write(bs)
		m.writer.Close()
	}
	return
}

func (m *UPRManager) Proposed(u string,
	p int) (pd *tesis.PageD, e error) {
	if m.steSys != nil && m.steSys.UsrAct != nil &&
		m.steSys.UsrAct[u] != nil {
		var a, b []interface{}
		a, pd = ConvDiffI(m.steSys.UsrAct[u].Proposed),
			&tesis.PageD{PageN: p}
		b, _, _, pd.Total = pageSlice(a, 10, p)
		pd.DiffP = ConvIDiff(b)
	}
	return
}

func ConvDiffI(ds []tesis.Diff) (r []interface{}) {
	r = make([]interface{}, len(ds))
	for i, j := range ds {
		r[i] = j
	}
	return
}

func ConvIDiff(is []interface{}) (r []tesis.Diff) {
	r = make([]tesis.Diff, len(is))
	for i, j := range is {
		r[i] = j.(tesis.Diff)
	}
	return
}

func (m *UPRManager) Pending(p int) (d *tesis.PageD, e error) {
	var r []tesis.Diff
	r = m.steSys.Pending //FIXME m.usrDt[u] = nil
	if r != nil {
		var t []tesis.Diff
		var a, b []interface{}
		a = ConvDiffI(r)
		var ps int
		b, _, _, ps = pageSlice(a, m.pgLen, p)
		t = ConvIDiff(b)
		d = &tesis.PageD{Total: ps, PageN: p, DiffP: t}
	} else {
		m.steSys.Pending = make([]tesis.Diff, 0)
		d = new(tesis.PageD)
	}
	return
}

func (m *UPRManager) RevertProp(u string, r []string) (e error) {
	if m.steSys == nil || m.steSys.UsrAct == nil ||
		m.steSys.UsrAct[u] == nil {
		e = fmt.Errorf("User %s has no activity", u)
	} else {
		// { m.steSys.UsrAct[u].Proposed ≠ nil }
		var rd []tesis.Diff
		rd = tesis.CreateDiff(r)
		var a, b, c, e []tesis.Eq
		a, b = tesis.ConvDiffEq(m.steSys.UsrAct[u].Proposed),
			tesis.ConvDiffEq(rd)
		c, e = tesis.DiffInt(a, b)
		m.steSys.UsrAct[u].Proposed, m.steSys.Pending =
			tesis.ConvEqDiff(c),
			append(m.steSys.Pending,
				tesis.ConvEqDiff(e)...)
	}
	return
}

func pageSlice(s []interface{}, n, p int) (t []interface{}, a, b, ps int) {
	// { n ≠ 0 }
	var rm int
	ps = len(s) / n //amount of pages
	rm = len(s) % n //amount of elements in reminder page
	t = make([]interface{}, 0)
	if rm != 0 {
		// {there is reminder page}
		ps++
		// {ps is the amount of pages,
		//  including the case where there is reminder}
	}
	if 0 <= p && p < ps {
		if p == ps-1 {
			// { p = ps-1 ∧ (rm = 0  ∨  rm ≠ 0) }
			// rm = 0  ⇒  len.s = p*(n+1)
			a, b = p*n, len(s)
			// { a = p*n ∧ b = len.s}
		} else {
			a, b = p*n, (p+1)*n
			// { a = p*n ∧ b = (p+1)*n }
		}
		// { r[a:b] = page.r.p }
		t = s[a:b]
		// { pg = page.r.p ≡ 0 ≤ p < len.(page.r) }
	}
	return
}
