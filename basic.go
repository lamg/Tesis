package syncbd

func Applicar(o Originator, r Addetor) {
	// { P.o.r }
	for o.HaActual() {
		// { P.o.r ∧ o.HaProxime() }
		r.Adder(o.Actual())
		// { Q.o.r }
		o.Proxime()
		// { P.o.r
	}
	// { P.o.r ∧ ¬o.HaProxime() }
}
