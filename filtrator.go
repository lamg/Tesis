package syncbd

type Filtrator interface {
	Originator
	Addetor
	Predicator
}

type SFilt struct {
	o Originator
	a Addetor
	p Predicator
}

func NewSFilt(o Originator, a Addetor, p Predicator) (r *SFilt) {
	r = &SFilt{o, a, p}
	return
}

func (s *SFilt) Adder(x interface{}) {
	if s.p.Ver(x) {
		s.a.Adder(x)
	}
}

func (s *SFilt) HaActual() (b bool) {
	b = s.o.HaActual()
	return
}

func (s *SFilt) Actual() (r interface{}) {
	r = s.o.Actual()
	return
}

func (s *SFilt) Proxime() {
	s.o.Proxime()
}

func (s *SCerc) Remontar() {
	s.o.Remontar()
}
