package tesis

import (
	"time"
	"github.com/dgrijalva/jwt-go"
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
	user string
	pass string
}

type Portal interface {
	Auth(Credentials) (jwt.Token, error)
}
