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
	var u, p, fname, rAdr, sf *string
	u, p, fname, rAdr, sf =
		flag.String("u", "", "User"),
		flag.String("p", "", "Password"),
		flag.String("f", "", "StateSystem file"),
		flag.String("rAdr", "", "Receiver address"),
		flag.String("sf", "", "Suffix of AD receiver")

	flag.Parse()
	if *rAdr == "" {
		e = tesis.CmbE(e, "Not supplied receptor address")
	}
	if *fname == "" {
		e = tesis.CmbE(e, "StateSystem file not defined")
	}
	if *p == "" {
		e = tesis.CmbE(e, "Password not defined")
	}
	if *u == "" {
		e = tesis.CmbE(e, "User not defined")
	}
	var fl io.ReadWriteCloser
	if e == nil {
		// { UserStr.u ∧ UserPass.p ∧ Filename.fname }
		fl, e = tesis.NewFileHandler(*fname)
		// { UserStr.u ∧ UserPass.p ∧ io.ReaderCloser.fl
		//   ≢ e = nil }
	}
	var bs []byte
	if e == nil {
		// { UserStr.u ∧ UserPass.p ∧ io.ReaderCloser.fl }
		bs, e = ioutil.ReadAll(fl)
		// { UserStr.u ∧ UserPass.p ∧ contents.fl = bs
		//   ≢ e = nil }
	}
	var ss *tesis.StateSys
	if e == nil {
		// { UserStr.u ∧ UserPass.p ∧ contents.fl = bs }
		ss = new(tesis.StateSys)
		e = json.Unmarshal(bs, ss)
		// { UserStr.u ∧ UserPass.p ∧ StateSys.ss ≢ e = nil }
	}
	var rcp tesis.RecordReceptor
	if e == nil {
		var l *db.LDAPAuth
		// { UserStr.u ∧ UserPass.p ∧ Address.rAdr ∧ Suff.sf }
		l, e = db.NewLDAPAuth(*rAdr, *sf)
		var b bool
		if e == nil {
			b, e = l.Authenticate(*u, *p)
		}
		if b {
			rcp = l
		} else {
			e = fmt.Errorf("Failed authentication")
		}
		// { UserStr.u ∧ RecordReceptor.rcp ∧
		//   StateSys.ss ≢ e = nil }
	}
	println("ok")
	if e == nil && ss.UsrAct != nil &&
		ss.UsrAct[*u] != nil &&
		ss.UsrAct[*u].Proposed != nil &&
		len(ss.UsrAct[*u].Proposed) != 0 {
		// { RecordReceptor.rcp ∧ UserStr.u ∧ StateSys.ss }
		var pr tesis.Reporter
		pr = tesis.NewPRpr()
		e = ss.SyncPend(rcp, *u, pr)
		// { written pending diffs to AD }
	} else if e == nil {
		log.Printf("No proposed changes for %s", *u)
	}
	// { written pending diffs to AD ≢ e = nil }
	var rs []byte
	if e == nil {
		rs, e = json.MarshalIndent(ss, "", "\t")
	}
	if e == nil {
		_, e = fl.Write(rs)
	}
	if e == nil {
		fl.Close()
	}
	// { written changes to ss to file ≢ e = nil }
	if e != nil {
		log.Fatal(e)
	}
}
