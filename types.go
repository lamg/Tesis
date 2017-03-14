package tesis

import (
	"github.com/dgrijalva/jwt-go"
)

type Info struct {
	Name string
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
	//bitacora de cambios hechos por el usuario
	UserInfo(string) (*Info, error)
}

type DummyAuth struct {
}

func (d *DummyAuth) Authenticate(u, p string) (b bool) {
	b = u == p
	return
}

type DummyManager struct {
}

func (m *DummyManager) UserInfo(u string) (inf *Info, e error) {
	inf = &Info{Name: u}
	return
}

type AccId string   //account id
type SStatus string //synchronization status

type AccMatch struct {
	DBId        AccId
	ADId        string
	MissingID   string
	MissingName string
	SrcDB       string
}

type Synchronizer interface {
	//synchronize a list of accounts
	Synchronize([]AccMatch) error

	//get the candidates for synchronization (who and why?)
	Candidates() ([]AccMatch, error)
}

// isCandidate ≡ ¬hasId ∨ existsSimilar
// existsSimilar ≡ toLowerEq ∨ unAccentEq

type DummySync struct {
	synced bool
}

func NewDummySync() (s *DummySync) {
	s = &DummySync{synced: false}
	return
}

func (s *DummySync) Candidates() (a []AccMatch, e error) {
	if !s.synced {
		a = []AccMatch{
			AccMatch{DBId: "0", ADId: "Coco", MissingID: "8901191122"},
			AccMatch{DBId: "1", MissingName: "Luis"},
		}
	}
	//iterate DB and filter comparing with AD
	return
}

func (s *DummySync) Synchronize(a []AccMatch) (e error) {
	s.synced = true
	return
}
