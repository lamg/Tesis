package db

import (
	"encoding/json"
	"github.com/lamg/tesis"
	"io"
	"io/ioutil"
)

type UPRManager struct {
	usrDB  tesis.UserDB
	steSys *tesis.StateSys
	// { JSONRepr.usrDt = contents.dtFile
	//   ≡ exists.dtFile }
	writer io.Writer
	pgLen  int // { pgLen ≥ 0 }
}

func NewUPRManager(f io.Reader, w io.Writer, a tesis.UserDB) (m *UPRManager, e error) {
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
			writer: w,
			pgLen:  10,
		}
	}
	return
}

func (m *UPRManager) Authenticate(u, p string) (b bool, e error) {
	b, e = m.usrDB.Authenticate(u, p)
	return
}

func (m *UPRManager) UserInfo(u string) (n *tesis.UserInfo, e error) {
	n, e = m.usrDB.UserInfo(u)
	return
}

func (m *UPRManager) Record(u string, p int) (c *tesis.PageC, e error) {
	var r []tesis.Change
	r = m.steSys.UsrAct[u].Record
	if r != nil {
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
	} else {
		m.steSys.UsrAct[u] = new(tesis.Activity)
		c = new(tesis.PageC)
	}
	return
}

func (m *UPRManager) Propose(u string, d []tesis.Diff) (e error) {
	if m.steSys == nil {
		m.steSys = new(tesis.StateSys)
	}
	if m.steSys.UsrAct == nil {
		m.steSys.UsrAct = make(map[string]*tesis.Activity)
	}
	var r *tesis.Activity
	r = m.steSys.UsrAct[u]
	if r == nil {
		r = new(tesis.Activity)
		m.steSys.UsrAct[u] = r
	}

	if r.Proposed == nil {
		r.Proposed = make([]tesis.Diff, 0)
	}
	var f, g, h, l []tesis.Eq
	f, g = make([]tesis.Eq, len(m.steSys.Pending)), make([]tesis.Eq, len(d))
	for i, j := range m.steSys.Pending {
		f[i] = j
	}
	for i, j := range d {
		g[i] = j
	}
	h, l = tesis.DiffInt(f, g)
	var k, n []tesis.Diff
	k, n = make([]tesis.Diff, len(h)), make([]tesis.Diff, len(l))
	for i, j := range h {
		k[i] = j.(tesis.Diff)
	}
	for i, j := range l {
		n[i] = j.(tesis.Diff)
	}

	//{ k = 'k - d  ∧  n = 'k ∩ d }
	m.steSys.Pending = k
	r.Proposed = append(r.Proposed, n...)

	//
	var bs []byte
	if e == nil {
		bs, e = json.MarshalIndent(m.steSys, "", "\t")
	}
	if e == nil {
		_, e = m.writer.Write(bs)
	}
	// { contents.(m.dtFile) = JSONRep.(m.usrDt) ≡ e = nil }
	return
}

func (m *UPRManager) Pending(p int) (d *tesis.PageD, e error) {
	var r []tesis.Diff
	r = m.steSys.Pending //FIXME m.usrDt[u] = nil
	if r != nil {
		if r == nil {
			r = make([]tesis.Diff, 1)
			d = new(tesis.PageD)
		} else {
			var t []tesis.Diff
			var a, b []interface{}
			a = make([]interface{}, len(r))
			for i, j := range r {
				a[i] = j
			}
			var ps int
			b, _, _, ps = pageSlice(a, m.pgLen, p)
			t = make([]tesis.Diff, len(b))
			for i, j := range b {
				t[i] = j.(tesis.Diff)
			}
			d = &tesis.PageD{Total: ps, PageN: p, DiffP: t}
		}
	} else {
		d = new(tesis.PageD)
	}
	return
}

func pageSlice(s []interface{}, n, p int) (t []interface{}, a, b, ps int) {
	var rm int
	ps = len(s) / n //amount of pages
	rm = len(s) % n //amount of elements in reminder page
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
