package tesis

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Info struct {
	Name    string     `json:"name"`
	Record  []Change   `json:"changeLog"`
	Matches []AccMatch `json:"matches"`
}

type Change struct {
	Time time.Time  `json:"time"`
	SRec []AccMatch `json:"srec"`
	FRec []AccMatch `json:"frec"`
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
type AccMatch struct {
	DBId    AccId
	ADId    string
	SrcIN   string
	ADName  string
	SrcName string
	SrcDB   string
}

type DBRecord struct {
	Id   AccId  //database key field
	IN   string //identity number
	Name string //person name
}

type Synchronizer interface {
	//synchronize a list of accounts
	Synchronize(string, []AccMatch) error

	//get the candidates for synchronization (who and why?)
	Candidates() ([]AccMatch, error)
}

// isCandidate ≡ ¬hasId ∨ existsSimilar
// existsSimilar ≡ toLowerEq ∨ unAccentEq

func EqAccMatch(a, b *AccMatch) (r bool) {
	r = a.ADId == b.ADId && a.ADName == b.ADName &&
		a.DBId == b.DBId && a.SrcDB == b.SrcDB &&
		a.SrcIN == b.SrcIN && a.SrcName == b.SrcName
	return
}
