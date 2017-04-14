package tesis

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

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

type DBManager interface {
	Authenticate(user, password string) bool
	UserInfo(string) (*UserInfo, error)
	Record(string, int) (*PageC, error)
	Propose(string, []Diff) error
	Pending(string, int) (*PageD, error)
}

type Diff struct {
	LDAPRec, DBRec   DBRecord
	Src              string
	Exists, Mismatch bool
}

type DBRecord struct {
	Id   string //database key field
	IN   string //identity number
	Name string //person name
}

// isCandidate ≡ ¬hasId ∨ existsSimilar
// existsSimilar ≡ toLowerEq ∨ unAccentEq

func EqDBRecord(a, b *DBRecord) (r bool) {
	r = a.Id == b.Id
	return
}
