package tesis

// This algorithm is a "descendant" of DiffInt
// c = a - b
// d and e are the couples of similar not equal elements
// f = b - c
func DiffSym(a, b []Sim, rp Reporter) (c, d, e, f []Sim) {
	var i, j, k, l int //i,j for a and k,l for b
	var tot, prog float32
	i, j, k, l, c, d, e, f, tot, prog = 0, 0, 0, 0,
		make([]Sim, 0, len(a)),
		make([]Sim, 0, max(len(a), len(b))),
		make([]Sim, 0, max(len(a), len(b))),
		make([]Sim, 0, len(b)),
		float32(len(a)*len(b)),
		0
	for !(i == len(a) && k == len(b)) {
		var ra, rb bool
		prog = float32(i*k) / tot
		rp.Progress(prog)
		ra, rb = i != len(a) && j != len(b) &&
			a[i].Similar(b[j]),
			k != len(b) && l != len(a) && b[k].Similar(a[l])
		if ra || rb {
			if ra {
				// { a.i ∈ a ∩ b }
				// a.i and b.j are equal ∨ a.i and b.j are similar
				if !a[i].Equals(b[j]) {
					d, e = append(d, a[i]), append(e, b[j])
				}
				// a.i and b.j are equal ∨ a.i and b.j are
				// stored in correspondent indexes of d and e
				i = i + 1
			}
			if rb {
				// { b.k ∈ a ∩ b }
				// not d = append(d, b[k])
				// for avoiding repetition
				k = k + 1
			}
			// { a.i or b.k was found in the other array }
		} else if i != len(a) && j == len(b) ||
			k != len(b) && l == len(a) {
			if i != len(a) && j == len(b) {
				// { a.i ∈ a ∧ a.i ∉ b }
				c, i, j = append(c, a[i]), i+1, 0
			}
			if k != len(b) && l == len(a) {
				// { b.k ∈ b ∧ b.k ∉ a }
				f, k, l = append(f, b[k]), k+1, 0
			}
			// { a.i or b.k wasn't found in other array }
		} else if i != len(a) && j != len(b) ||
			k != len(b) && l != len(a) {
			// { a.i ≠ b.j ∨ b.k ≠ b.l  }
			if i != len(a) && j != len(b) {
				j = j + 1
			}
			if k != len(b) && l != len(a) {
				l = l + 1
			}
			// { the current element in the other array
			//   for comparing with a.i or b.k is discarded }
		}
		// { every element until i in a, and until k in b
		//   is classified or there is elements in the other
		//   array for comparing }
	}
	// { c = a - b ∧ d,e have similar (not equal) elements
	//   in homologal indexes ∧ e = b - a }
	return
}

func max(x, y int) (r int) {
	if x >= y {
		r = x
	} else {
		r = y
	}
	// { r = x ↑ y }
	return
}

func DiffInt(a, b []Eq) (c, e []Eq) {
	var i, j int
	var ok, d bool
	i, j, d, ok, c, e = 0, 0, true, false,
		make([]Eq, 0, len(a)),
		make([]Eq, 0, len(b))
	for d {
		if (i != len(a) && j == len(b)) || ok {
			if !ok {
				c = append(c, a[i])
			} else {
				e = append(e, a[i])
			}
			i, j, d, ok = i+1, 0, true, false
		} else if i != len(a) && j != len(b) && !ok {
			ok, j = a[i].Equals(b[j]), j+1
		} else if i == len(a) && !ok {
			d = false
		}
	}
	// { c = a - b ∧ e = a ∩ b }
	return
}

func ConvSim(s []DBRecord) (r []Sim) {
	r = make([]Sim, len(s))
	for i, j := range s {
		r[i] = j
	}
	return
}

func ConvDBR(s []Sim) (r []DBRecord) {
	r = make([]DBRecord, len(s))
	for i, j := range s {
		r[i] = j.(DBRecord)
	}
	return
}

/*
x ≡ i ≠ len(a)
y ≡ j ≠ len(b)
Calculating the negation of the last guard
  (x ∧ ¬y) ∨ ok ∨ (x ∧ y ∧ ¬ok)
≡ { complement  }
  (x ∧ ¬y) ∨ ok ∨ (x ∧ y)
≡ { ∧ over ∨ }
  (x ∧ (y ∨ ¬y)) ∨ ok
≡ { negation, unit }
  x ∨ ok
*/
