package main

import (
	"flag"
	"github.com/lamg/tesis"
	"github.com/lamg/tesis/db"
	"github.com/lamg/tesis/http"
	"log"
	"os"
	"os/signal"
)

func main() {
	var hp, cr, ky, lda, suf, dtf *string
	var dmy *bool

	hp, cr, ky, lda, suf, dtf, dmy =
		flag.String("p", ":10443", "Port to serve"),
		flag.String("c", "cert.pem", "PEM certificate file"),
		flag.String("k", "key.pem", "PEM key file"),
		flag.String("ls", "ad.upr.edu.cu:636",
			"LDAP server address"),
		flag.String("sf", "@upr.edu.cu", "Account suffix"),
		flag.String("df", "dtFile.json",
			"Activity record in JSON format"),
		flag.Bool("d", false,
			"Use dummy authentication instead LDAP")
	flag.Parse()
	var qr tesis.UserDB
	var e error
	var um *db.UPRManager

	if *dmy {
		qr = tesis.NewDummyManager()
	} else {
		qr, e = db.NewLDAPAuth(*lda, *suf)
	}
	var fl *tesis.FileHandler
	if e == nil {
		fl, e = tesis.NewFileHandler(*dtf)
	}
	if e == nil {
		um, e = db.NewUPRManager(fl, qr)
	}
	if e == nil {
		var sg chan os.Signal
		sg = make(chan os.Signal, 1)
		signal.Notify(sg, os.Interrupt)
		go closeM(sg, um)
		http.ListenAndServe(*hp, um, *cr, *ky)
	}
	if e != nil {
		log.Fatal(e)
	}
}

func closeM(sg chan os.Signal, m tesis.DBManager) {
	<-sg
	m.Close()
	os.Exit(0)
}
