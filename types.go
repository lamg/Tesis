package tesis

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Info struct {
	UsrInf   *UserInfo `json:"userInfo"`
	Record   []Change  `json:"changeLog"`
	Proposed []Diff    `json:"proposed"`
	Pending  []Diff    `json:"pending"`
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

type Authenticator interface {
	Authenticate(user, password string) bool
}

type DBManager interface {
	UserInfo(string) (*Info, error)
	Synchronizer
}

type AccId string   //account id
type SStatus string //synchronization status

// AccMatch ≡ MissingIDMatch ∨ DiffNamesMatch
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

type Synchronizer interface {
	//synchronize a list of accounts
	Synchronize(string, []Diff) error

	//get the candidates for synchronization (who and why?)
	Candidates() ([]Diff, error)
}

// isCandidate ≡ ¬hasId ∨ existsSimilar
// existsSimilar ≡ toLowerEq ∨ unAccentEq

func EqDBRecord(a, b *DBRecord) (r bool) {
	r = a.Id == b.Id
	return
}
