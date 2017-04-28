package tesis

import (
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"
	"unicode"
)

// This interface is an abstract program specification.
// The rest of types and procedures are defined for
// using them inside it, or its implementation.
type DBManager interface {
	Authenticate(user, password string) (bool, error)
	UserInfo(string) (*UserInfo, error)
	Record(string, int) (*PageC, error)
	Propose(string, []Diff) error
	Pending(int) (*PageD, error)
}

type RecordProvider interface {
	Records() ([]DBRecord, error)
	Name() string
}

type Activity struct {
	Record   []Change `json:"record"`
	Proposed []Diff   `json:"proposed"`
}

type StateSys struct {
	Pending []Diff               `json:"pending"`
	UsrAct  map[string]*Activity `json:"usrAct"`
}

type UserDB interface {
	Authenticate(string, string) (bool, error)
	UserInfo(string) (*UserInfo, error)
}

//TODO construct adequate interfaces

type PageN struct {
	PageN int `json:"pageN"`
}

type PageC struct {
	Total   int      `json:"total"`
	PageN   int      `json:"pageN"`
	ChangeP []Change `json:"changeP"`
}

type PageD struct {
	Total int    `json:"total"`
	PageN int    `json:"pageN"`
	DiffP []Diff `json:"diffP"`
}

type UserInfo struct {
	Name string `json:"name"`
}

type Error struct {
	Message string `json:"error"`
}

type Change struct {
	Time time.Time `json:"time"`
	SRec []Diff    `json:"srec"`
	FRec []Diff    `json:"frec"`
}

type Credentials struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type User struct {
	UserName string `json:"user"`
	jwt.StandardClaims
}

type Diff struct {
	LDAPRec  DBRecord `json:"ldapRec"`
	DBRec    DBRecord `json:"dbRec"`
	Src      string   `json:"src"`
	Exists   bool     `json:"exists"`
	Mismatch bool     `json:"mismatch"`
}

type DBRecord struct {
	Id string `json:"id"`
	//database key field
	IN string `json:"in"`
	//identity number
	Name string `json:"name"`
	//person name
}

// isCandidate ≡ ¬hasId ∨ existsSimilar
// existsSimilar ≡ toLowerEq ∨ unAccentEq

type Eq interface {
	Equals(interface{}) bool
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
	t = strings.Replace(strings.ToLower(s), " ", "", -1)
	var rd io.Reader
	rd = strings.NewReader(t)
	var isMn func(rune) bool
	isMn = func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}
	var tr transform.Transformer
	tr = transform.Chain(norm.NFD,
		transform.RemoveFunc(isMn), norm.NFC)
	rd = transform.NewReader(rd, tr)
	var bs []byte
	bs, _ = ioutil.ReadAll(rd)
	t = string(bs)
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

type Sim interface {
	Similar(interface{}) bool
	Eq
}

type Nat int

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

// This algorithm is a "descendant" of DiffInt
// c = a - b
// d and e are the couples of similar elements
// f = b - c
func DiffSym(a, b []Sim) (c, d, e, f []Sim) {
	var i, j, k, l int //i,j for a and k,l for b
	i, j, k, l, c, d, e, f = 0, 0, 0, 0,
		make([]Sim, 0, len(a)),
		make([]Sim, 0, max(len(a), len(b))),
		make([]Sim, 0, max(len(a), len(b))),
		make([]Sim, 0, len(b))
	for !(i == len(a) && k == len(b)) {
		var ra, rb bool
		ra, rb = i != len(a) && j != len(b) &&
			a[i].Similar(b[j]),
			k != len(b) && l != len(a) && b[k].Similar(a[l])
		if ra || rb {
			if ra {
				log.Print(i)
				// { a.i ∈ a ∩ b }
				// a.i and b.j are equal ∨ a.i and b.j are similar
				d, e = append(d, a[i]), append(e, b[j])
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
