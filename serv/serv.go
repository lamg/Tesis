package main

import (
	"flag"
	"github.com/lamg/tesis"
	"github.com/lamg/tesis/db"
	"github.com/lamg/tesis/http"
	"log"
)

func main() {
	var hp, cr, ky, lds, suf, dtf *string
	var ldp *int
	var dmy *bool

	hp, cr, ky, lds, suf, dtf, ldp, dmy =
		flag.String("p", ":10443", "Port to serve"),
		flag.String("c", "cert.pem", "PEM certificate file"),
		flag.String("k", "key.pem", "PEM key file"),
		flag.String("ls", "ad.upr.edu.cu",
			"LDAP server address"),
		flag.String("sf", "@upr.edu.cu", "Account suffix"),
		flag.String("df", "dtFile.json",
			"Activity record in JSON format"),
		flag.Int("lp", 636, "LDAP server port"),
		flag.Bool("d", false,
			"Use dummy authentication instead LDAP")
	flag.Parse()
	var qr tesis.UserDB
	var e error
	var um *db.UPRManager

	if *dmy {
		qr = tesis.NewDummyManager()
	} else {
		qr, e = db.NewLDAPAuth(*lds, *suf, *ldp)
	}

	if e == nil {
		qr = tesis.NewDummyManager()
		um, e = db.NewUPRManager(*dtf, qr)
	}
	if e == nil {
		http.ListenAndServe(*hp, um, *cr, *ky)
	}
	if e != nil {
		log.Fatal(e)
	}
}
