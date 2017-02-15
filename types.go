package tesis

import (
	"time"
	_ "github.com/dgrijalva/jwt-go"
)

//cuantos correos has mandado y recibido
//cuantos megas has consumido de tu cuenta
//de correo y de internet

type Info struct {
	SentMessages int
	RecvMessages int
	MailStorage  int
	InternetDwnl int
	WifiLogons   WifiL
}

type WifiL struct {
	Ip string
	Place string
	Date time.Time
}

type Credentials struct {
	User string
	Pass string
}

type Authenticator interface {
	Authenticate(user, password string) bool
}