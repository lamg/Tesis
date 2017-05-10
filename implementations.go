package tesis

import (
	"fmt"
	"os"
	"strings"
)

type Nat int //Implements Sim

func (n Nat) Equals(o interface{}) (b bool) {
	var m Nat
	m, b = o.(Nat)
	b = b && n == m
	return
}

func (n Nat) Similar(o interface{}) (b bool) {
	b = n.Equals(o)
	return
}

type TRpr struct {
	t   Logger
	Log bool
}

func NewTRpr(t Logger) (r *TRpr) {
	r = &TRpr{t: t}
	return
}

func (r *TRpr) Progress(p float32) {
	if r.Log {
		r.t.Logf("%.0f", p*100)
	}
}

type PRpr struct {
}

func NewPRpr() (r *PRpr) {
	r = new(PRpr)
	return
}

func (r *PRpr) Progress(p float32) {
	fmt.Fprintf(os.Stderr, "%.3f%s\r", p*100, "%")
}

func CmbE(e error, s string) (d error) {
	d = fmt.Errorf("%s ∧ %s", e.Error(), s)
	return
}

func (d DBRecord) Equals(c interface{}) (b bool) {
	var x DBRecord
	x, b = c.(DBRecord)
	b = b && x.Id == d.Id && x.IN == d.IN &&
		x.Name == d.Name
	return
}

func (d DBRecord) Similar(o interface{}) (b bool) {
	var e DBRecord
	e, b = o.(DBRecord)
	b = b && (toStd(d.Name) == toStd(e.Name) ||
		d.Equals(e))
	return
}

func toStd(s string) (t string) {
	t = strings.Map(func(x rune) (y rune) {
		if x == 'á' {
			y = 'a'
		} else if x == 'é' {
			y = 'e'
		} else if x == 'í' {
			y = 'i'
		} else if x == 'ó' {
			y = 'o'
		} else if x == 'ú' {
			y = 'u'
		} else if x == 'ñ' {
			y = 'n'
		} else if x == 'Á' {
			y = 'a'
		} else if x == 'É' {
			y = 'e'
		} else if x == 'Í' {
			y = 'i'
		} else if x == 'Ó' {
			y = 'o'
		} else if x == 'Ú' {
			y = 'u'
		} else if x == 'Ñ' {
			y = 'n'
		} else if x == ' ' {
			y = -1
		} else {
			y = x
		}
		return
	}, s)
	return
}

func (d Diff) Equals(c interface{}) (b bool) {
	var x Diff
	x, b = c.(Diff)
	b = b && d.DBRec.Equals(x.DBRec) &&
		d.Exists == x.Exists &&
		d.Mismatch == x.Mismatch &&
		d.Src == x.Src
	return
}
