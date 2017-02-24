package tesis

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

//cuantos correos has mandado y recibido
//cuantos megas has consumido de tu cuenta
//de correo y de internet

type Info struct {
	SentMessages int `json: "sentMessages"`
	RecvMessages int
	MailStorage  int
	InternetDwnl int
	WifiLogons   []WifiL
}

type WifiL struct {
	Ip    string
	Place string
	Date  time.Time
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
	inf = &Info{
		SentMessages: 18,
		RecvMessages: 40,
		MailStorage:  67,
		InternetDwnl: 87,
		WifiLogons: []WifiL{
			WifiL{
				Ip:    "192.168.0.10",
				Place: "Rectoría",
				Date:  time.Now(),
			},
		},
	}
	return
}
