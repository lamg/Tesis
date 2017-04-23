package tesis

import (
	"github.com/dgrijalva/jwt-go"
	"time"
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

func (d Diff) Equals(c interface{}) (b bool) {
	var x Diff
	x, b = c.(Diff)
	b = b && d.DBRec.Equals(x.DBRec) &&
		d.Exists == x.Exists &&
		d.Mismatch == x.Mismatch &&
		d.Src == x.Src
	return
}

type Nat int

func (n Nat) Equals(o interface{}) (b bool) {
	var m Nat
	m, b = o.(Nat)
	b = b && n == m
	return
}

func SymDiff(a, b []Eq) (c, e []Eq) {
	var i, j int
	var ok, d bool
	i, j, d, ok, c, e = 0, 0, true, false, make([]Eq, 0, len(a)), make([]Eq, 0, len(b))
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
	// {⟨∀i: i ∈ a ∧ i ∉ b: i ∈ c⟩ ∧
	//  ⟨∀i: i ∈ a ∧ i ∈ b: i ∈ e⟩ }
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
