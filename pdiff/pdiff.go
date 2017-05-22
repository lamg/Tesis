package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/lamg/tesis"
	"github.com/lamg/tesis/db"
	"io"
	"io/ioutil"
	"log"
)

func main() {
	var e error
	var usrDB, ldapAdr, sigenuAdr, assetAdr,
		us, ps, ua, pa *string //user database, password database,
	//user AD, pasword AD
	usrDB, ldapAdr, sigenuAdr, assetAdr, us, ps, ua, pa =
		flag.String("d", "dtFile.json", "JSON formated StateSys"),
		flag.String("la", "10.2.24.35:636", "LDAP address"),
		flag.String("sa", "10.2.24.117/sigenu", "SIGENU address"),
		flag.String("aa", "", "ASSET address"),
		flag.String("us", "", "User to access database"),
		flag.String("ps", "", "Password to access database"),
		flag.String("ua", "", "User to access AD"),
		flag.String("pa", "", "Password to access AD")

	flag.Parse()
	if *assetAdr != "" {
		e = fmt.Errorf("ASSET synchronization not implemented")
	}
	var lp, sg tesis.RecordProvider
	if e == nil {
		sg, e = db.NewPSProvider(*us, *ps, *sigenuAdr, -1)
	}
	if e == nil {
		lp, e = db.NewLDAPProv(*ua, *pa, *ldapAdr, -1)
	}
	var ds []tesis.Diff
	if e == nil {
		var pr *tesis.PRpr
		pr = tesis.NewPRpr()
		ds, e = db.PDiff(sg, lp, pr)
		// { ds contains inconsistent records in lp, according
		//   sg records ≢ e = nil }
	}
	var fl io.ReadWriteCloser
	if e == nil {
		fl, e = tesis.NewFileHandler(*usrDB)
	}
	var bs []byte
	if e == nil {
		bs, e = ioutil.ReadAll(fl)
	}
	var ss *tesis.StateSys
	if e == nil {
		ss = new(tesis.StateSys)
		json.Unmarshal(bs, ss)
	}
	var rs []byte
	if e == nil {
		ss.Pending = ds
		rs, e = json.MarshalIndent(ss, "", "\t")
	}
	if e == nil {
		_, e = fl.Write(rs)
	}
	if e == nil {
		fl.Close()
	}
	if e != nil {
		log.Fatal(e)
	}
}
