package main

import (
	"fmt"
	"github.com/lamg/tesis"
	"github.com/lamg/tesis/http"
	"os"
)

func main() {
	var hp string
	var au tesis.Authenticator
	var qr tesis.DBManager
	var e error

	hp = "localhost:10443"
	// au, e = tesis.NewLDAPAuth("ad.upr.edu.cu", "@upr.edu.cu", 636)
	au = new(tesis.DummyAuth)
	if e == nil {
		qr = tesis.NewDummyManager()
		//{ cwd contains files used by server }
		http.ListenAndServe(hp, au, qr)
	}
	if e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
	}
}
