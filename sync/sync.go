package main

import (
	"encoding/json"
	"flag"
	"github.com/lamg/tesis"
	"github.com/lamg/tesis/db"
	"io"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	var e error
	var usrDB, ldapAdr, sigenuAdr, assetAdr,
		user, pass *string
	usrDB, ldapAdr, sigenuAdr, assetAdr, user, pass =
		flag.String("d", "dtFile.json", "JSON formated StateSys"),
		flag.String("la", "10.2.24.35:636", "LDAP address"),
		flag.String("pa", "10.2.24.117/sigenu", "SIGENU address"),
		flag.String("aa", "", "ASSET address"),
		flag.String("us", "", "User to access databases"),
		flag.String("ps", "", "Password to access databases")
	flag.Parse()
	println(assetAdr)
	var lp, sg tesis.RecordProvider
	sg, e = db.NewPSProvider("lamg", "hqmnv78", *sigenuAdr, -1)
	if e == nil {
		lp, e = db.NewLDAPProv(*user, *pass, *ldapAdr, -1)
	}
	var ds []tesis.Diff
	if e == nil {
		var pr *tesis.PRpr
		pr = tesis.NewPRpr()
		ds, e = db.Sync(sg, lp, pr)
	}
	var fl io.ReadCloser
	if e == nil {
		fl, e = os.Open(*usrDB)
	}
	var bs []byte
	if e == nil {
		bs, e = ioutil.ReadAll(fl)
		fl.Close()
	}
	var ss *tesis.StateSys
	if e == nil {
		ss = new(tesis.StateSys)
		e = json.Unmarshal(bs, ss)
	}
	var rs []byte
	if e == nil {
		ss.Pending = ds
		rs, e = json.Marshal(ss)
	}
	var f io.WriteCloser
	if e == nil {
		f, e = os.Create(*usrDB)
	}
	if e == nil {
		_, e = f.Write(rs)
	}
	if e == nil {
		f.Close()
	}
	if e != nil {
		log.Fatal(e)
	}
}
