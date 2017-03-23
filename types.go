package tesis

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Info struct {
	Name    string     `json:"name"`
	Error   string     `json:"error"`
	Record  []Change   `json:"record"`
	Matches []AccMatch `json:"matches"`
}

type Change struct {
	Time *time.Time
	SRec []AccMatch
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
	//todos los usuarios reciben la misma informacion
	//del estado del sistema?
	UserInfo(string) (*Info, error)
	Synchronizer
}

type AccId string   //account id
type SStatus string //synchronization status

// AccMatch ≡ MissingIDMatch ∨ DiffNamesMatch
type AccMatch struct {
	DBId    AccId
	ADId    string
	SrcID   string
	ADName  string
	SrcName string
	SrcDB   string
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
		a.SrcID == b.SrcID && a.SrcName == b.SrcName
	return
}
