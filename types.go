package tesis

import (
	"github.com/dgrijalva/jwt-go"
)

type Info struct {
	name string
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
	inf = &Info{name: u}
	return
}
