package syncbd

type Cercator interface {
	Originator
	Addetor
}

type SCerc struct {
	o Originator
	a Addetor
	p Predicator
	b bool
}

func NewSCerc(o Originator, a Addetor,
	p Predicator) (c *SCerc) {
	c = &SCerc{o, a, p}
	return
}

func (c *SCerc) Adder(x interface{}) {
	if c.p.Ver(x) {
		c.a.Adder(x)
		b = false
	}
}

func (c *SCerc) HaActual() (b bool) {
	b = c.b && c.o.HaActual()
	return
}

func (c *SCerc) Actual() (r interface{}) {
	r = c.o.Actual()
	return
}

func (c *SCerc) Proxime() {
	if c.b {
		c.o.Proxime()
	}
}

func (c *SCerc) Remontar() {
	c.o.Remontar()
}
