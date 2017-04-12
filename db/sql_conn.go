package db

import (
	"database/sql"
	"github.com/lamg/tesis"
)

type DBSync struct {
	l *LDAPAuth
}

func NewDBSync(l *LDAPAuth) (d *DBSync, e error) {
	d.l = l
	return
}

func (d *DBSync) Synchronize(a []tesis.AccMatch) (e error) {
	return
}

func (d *DBSync) Candidates() (a []tesis.AccMatch, e error) {
	var ldpR, dbR []tesis.DBRecord
	// all ldap students
	ldpR, e = d.l.GetUsers()
	if e == nil {
		dbR, e = AllStudents()
	}
	// find differences
	if e == nil {
		for _, i := range ldpR {
			for _, j := range dbR {

			}
			//add if i not found or inconsistent
		}
	}
	return
}

//all Candidates not synchronized are stored for
//future synchronizations
