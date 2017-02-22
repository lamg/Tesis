package main

import (
	"github.com/lamg/tesis"
)

func main() {
	var e error
	hp := "localhost:10443"
	ce := "cert.pem"
	ke := "key.pem"
	au := &DummyAuth{}
	qr := &DummyManager{}
	tesis.NewHTTPPortal(hp, ce, ke, au, qr)
}
