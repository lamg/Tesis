package syncbd

import (
	"strings"
)

func NewUnitorRegistros() (r Unitor) {
	r = &UReg{}
	return
}

type UReg struct {
}

func (u *UReg) Unir(a, b interface{}) (c interface{}) {
	var ak, bk, ok bool
	var r, s *Registro
	r, ak = a.(*Registro)
	s, bk = b.(*Registro)
	if ak && bk {
		ok = Similar(r, s)
	}
	if a || b {
		c = &Diff{Anterior: r, Posterior: s, Similar: ok}
	}
	return
}

func Similar(a, b *Registro) (c bool) {
	c = toStd(a.Nomine) == toStd(b.Nomine)
	return
}

type Registro struct {
	//Identitate in base de datos
	BDId   string
	Nomine string
	//Numero de identitate
	NumId     string
	Direction string
	Telephono string
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
