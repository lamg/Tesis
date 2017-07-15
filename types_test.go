package tesis

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDiffSym(t *testing.T) {
	//this proof has no sense for naturals
	var db0, db1, db2, db3, db4 DBRecord
	db0, db1, db2, db3, db4 = DBRecord{
		Id:   "516648e2:14fd74554bc:-2309",
		IN:   "95041823862",
		Name: "Yansel Acosta Sarabia",
		Addr: "Ave del Valle",
		Tel:  "448416",
	}, DBRecord{
		Name: "Yansel Acosta Sarabia",
	}, DBRecord{
		Id:   "-3457f221:13292017856:-6c8e",
		IN:   "91050527469",
		Name: "Yasniel Salabarría Castellanos",
		Addr: "Edif. E 14 Apto 40. López Peña",
		Tel:  "565458",
	}, DBRecord{
		IN:   "91050527469",
		Name: "Yasniel Salabarria Castellanos",
		Addr: "Edif. E 14 Apto 40. Lopez Pena, Artemisa, Cuba",
		Tel:  "565458",
	}, DBRecord{
		Name: "Coco",
	}

	var v, w, x, y, z []DBRecord
	v, w, x, y, z = []DBRecord{
		db0,
		db2,
		db4,
	},
		[]DBRecord{
			db1,
			db3,
		},
		[]DBRecord{
			db4,
		}, []DBRecord{
			db0,
			db2,
		}, []DBRecord{
			db1,
			db3,
		}
	// x = v - w ∧ y = v ∩ w ∧ z = w - x
	var a, b, c, d, e, f []Sim
	a, b = ConvSim(v), ConvSim(w)
	var rp *TRpr
	rp = NewTRpr(t)
	c, d, e, f = DiffSym(a, b, rp)
	require.True(t, eqsim(x, c))
	require.True(t, eqsim(y, d))
	require.True(t, eqsim(y, e))
	// error if but not significative
	// since PDiff ignores c and
	t.Log(len(f))
	for _, j := range f {
		t.Logf("%v", j)
	}
	require.True(t, eqsim(z, f))
}

func conv(x []Nat) (y []Sim) {
	y = make([]Sim, len(x))
	for i, j := range x {
		y[i] = j
	}
	return
}

func eqsim(x []DBRecord, y []Sim) (b bool) {
	var i int
	i = 0
	if len(x) == len(y) {
		for i != len(x) && x[i].Similar(y[i]) {
			i = i + 1
		}
	}
	b = i == len(x)
	return
}

func TestToStd(t *testing.T) {
	var p, q, r string
	p, q = "Luis Ángel Méndez Gort Ññ", "LuisangelMendezGortnn"
	r = toStd(p)
	require.EqualValues(t, q, r)
}
