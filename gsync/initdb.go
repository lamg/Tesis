package main

import (
	"encoding/json"
	"github.com/lamg/tesis"
	"github.com/lamg/tesis/db"
	"io"
	"io/ioutil"
	"os"
	"path"
)

const (
	//Configuration file name
	cfgFle = ".config/gsync/gsync.json"
	dtFile = ".local/share/gsync/dtFile.json"
)

type Cfg struct {
	//Path of system state persistence file
	DtFile string `json:"dtFile"`
	//Address and SSL port of Active Directory server
	ADAddr string `json:"adAddr"`
	//Address and SSL port of SIGENU server
	SGAddr string `json:"sgAddr"`
}

func initDB() (stSys tesis.DBManager, e error) {
	var cf, dt string
	if e == nil {
		var hm string
		hm = os.Getenv("HOME")
		cf, dt = path.Join(hm, cfgFle), path.Join(hm, dtFile)
	}
	var bs []byte
	if e == nil {
		bs, e = ioutil.ReadFile(cf)
		if e != nil {
			e = initPrs(cf, dt)
		}
	}
	var c *Cfg
	if e == nil {
		c = new(Cfg)
		e = json.Unmarshal(bs, c)
	}
	var rwc io.ReadWriteCloser
	if e == nil {
		rwc, e = tesis.NewFileHandler(c.DtFile)
	}
	var usDB tesis.UserDB
	if e == nil {
		usDB, e = db.NewLDAPAuth(c.ADAddr, "@upr.edu.cu")
	}
	if e == nil {
		stSys, e = db.NewUPRManager(rwc, usDB)
	}
	return
}

func initPrs(cf, dt string) (e error) {
	var ada, sga string
	var c *Cfg
	//TODO
	ada, sga = "10.2.24.35:636", "10.2.24.117:POSTGRE_PORT/sigenu"
	if e == nil {
		os.Mkdir(path.Dir(cf), os.ModeDir)
	}
	if e == nil {
		c = &Cfg{DtFile: dt, ADAddr: ada, SGAddr: sga}
	}
	var bs []byte
	if e == nil {
		bs, e = json.Marshal(c)
	}
	if e == nil {
		e = ioutil.WriteFile(cf, bs, 0666)
	}
	// { config file written in cf }

	if e == nil {
		os.Mkdir(path.Dir(dt), os.ModeDir)
	}
	if e == nil {
		var ste *tesis.StateSys
		ste = &tesis.StateSys{
			Pending: make([]tesis.Diff, 0),
			UsrAct:  make(map[string]*tesis.Activity),
		}
		bs, e = json.Marshal(ste)
	}
	if e == nil {
		e = ioutil.WriteFile(dt, bs, 0666)
	}
	// { steSys file written in dtf }
	return
}
