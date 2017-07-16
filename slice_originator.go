package syncbd

func NewSliceOriginator(r []interface{}) (o Originator) {
	o = &SLOrg{slice: r, indice: 0}
	return
}

type SLOrg struct {
	slice  []interface{}
	indice int
}

func (s *SLOrg) HaProxime() (r bool) {
	r = s.indice != len(s.slice)
	return
}

func (s *SLOrg) Proxime() (r interface{}) {
	r, s.indice = s.slice[s.indice], s.indice+1
	return
}

func (s *SLOrg) Remontar() {
	s.indice = 0
}
