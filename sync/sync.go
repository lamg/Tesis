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
	var u, p, fname, rAdr *string
	u, p, fname, rAdr =
		flag.String("u", "", "User"),
		flag.String("p", "", "Password"),
		flag.String("f", "", "StateSystem file"),
		flag.String("rAdr", "", "Receiver address")

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
	var fl io.ReadCloser
	if e == nil {
		// { UserStr.u ∧ UserPass.p ∧ Filename.fname }
		fl, e = os.Open(*fname)
		// { UserStr.u ∧ UserPass.p ∧ io.ReaderCloser.fl
		//   ≢ e = nil }
	}
	var bs []byte
	if e == nil {
		// { UserStr.u ∧ UserPass.p ∧ io.ReaderCloser.fl }
		bs, e = ioutil.ReadAll(fl)
		fl.Close()
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
		// { UserStr.u ∧ UserPass.p ∧ Address.rAdr }
		rcp, e = db.NewLDAPRecp(*rAdr, *u, *p)
		// { UserStr.u ∧ RecordReceptor.rcp ∧
		//   StateSys.ss ≢ e = nil }
	}
	if e == nil {
		// { RecordReceptor.rcp ∧ UserStr.u ∧ StateSys.ss }
		e = ss.SyncPend(rcp, *u)
		// { written pending diffs to AD }
	}
	// { written pending diffs to AD ≢ e = nil }
	if e != nil {
		log.Fatal(e)
	}
}
