package tesis

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Info struct {
	Name    string
	Record  []Change
	Matches []AccMatch
}

type Change struct {
	Time *time.Time
	SRec []AccMatch
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
	Synchronize([]AccMatch) error

	//get the candidates for synchronization (who and why?)
	Candidates() ([]AccMatch, error)
}

// isCandidate ≡ ¬hasId ∨ existsSimilar
// existsSimilar ≡ toLowerEq ∨ unAccentEq
