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
	var fs *http.ServFS

	hp, fs = "localhost:10443", &http.ServFS{"st", "cert.pem", "key.pem", []string{"st/index.html", "st/dash.html"}}

	// au, e = tesis.NewLDAPAuth("ad.upr.edu.cu", "@upr.edu.cu", 636)
	au = new(tesis.DummyAuth)
	if e == nil {
		qr = tesis.NewDummyManager()
		//{ cwd contains files used by server }
		http.ListenAndServe(hp, au, qr, fs)
	}
	if e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
	}
}
