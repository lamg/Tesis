package syncbd

type Cruciator interface {
	Originator
	Addetor
}

type SCruc struct {
	a, b Originator
	u    Unitor
	i, j int
	ant  bool
	act  interface{}
}

func NewAddetorCruciar(a, b Originator, c Addetor,
	u Unitor) (c *AddetorCruciar) {
	c = &AddetorCruciar{
		a: a, b: b,
		i: 0, j: 0,
		u:   u,
		ant: AleatoriBool(),
	}
}

func (ct *AddetorCruciar) HaActual() (b bool) {
	b = ct.a.HaActual() || ct.b.HaActual()
	return
}

func (ct *AddetorCruciar) Actual() (r interface{}) {
	// { ct.a.HaActual() ∨ ct.b.HaActual() }
	// TODO
	var a, b interface{}
	var o, p Originator
	if ct.ant {
		o, p = ct.a, ct.b
	} else {
		o, p = ct.b, ct.a
	}
	r = u.Unir(o.Actual(), p.Actual())
}
