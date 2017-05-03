package tesis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiffSym(t *testing.T) {
	var v, w, x, y, z []Nat
	v, w, x, y, z = []Nat{1, 2, 3, 4, 5}, []Nat{2, 4, 18},
		[]Nat{1, 3, 5}, []Nat{2, 4}, []Nat{18}
	// x = v - w ∧ y = v ∩ w ∧ z = w - x
	var a, b, c, d, e, f []Sim
	a, b = conv(v), conv(w)
	var rp *tRpr
	rp = NewTRpr(t)
	c, d, e, f = DiffSym(a, b, rp)
	assert.True(t, eqnat(x, c) &&
		eqnat(y, d) && eqnat(y, e) && eqnat(z, f))
}

func conv(x []Nat) (y []Sim) {
	y = make([]Sim, len(x))
	for i, j := range x {
		y[i] = j
	}
	return
}

func eqnat(x []Nat, y []Sim) (b bool) {
	var i int
	i = 0
	if len(x) == len(y) {
		for i != len(x) && x[i].Equals(y[i]) {
			i = i + 1
		}
	}
	b = i == len(x)
	return
}

func TestToStd(t *testing.T) {
	var p, q, r string
	p, q = "Luis Ángel Méndez Gort Ññ", "luisangelmendezgortnn"
	r = toStd(p)
	if !assert.EqualValues(t, q, r) {
		t.Log(r)
	}
}
