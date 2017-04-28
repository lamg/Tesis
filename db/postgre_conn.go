package db

import (
	"database/sql"
	"fmt"
	"github.com/lamg/tesis"
	_ "github.com/lib/pq"
)

func NewPSProvider(u, p string) (r tesis.RecordProvider, e error) {
	var l *PSProv
	l = new(PSProv)
	var s string
	s = fmt.Sprintf("postgres://%s:%s@10.2.24.117/sigenu", u, p)
	l.db, e = sql.Open("postgres", s)
	if e == nil {
		r = l
	}
	return
}

type PSProv struct {
	db *sql.DB
}

func (p *PSProv) Records() (s []tesis.DBRecord, e error) {
	var r *sql.Rows
	r, e = p.db.Query("SELECT id_student,identification,name,middle_name,last_name FROM student")
	s = make([]tesis.DBRecord, 0)
	for e == nil && r.Next() {
		var st tesis.DBRecord
		st = tesis.DBRecord{}
		var name, middle_name, last_name string
		e = r.Scan(&st.Id, &st.IN, &name, &middle_name, &last_name)
		if e == nil {
			st.Name = name + middle_name + last_name
			s = append(s, st)
		}
	}
	return
}

func (p *PSProv) Name() (s string) {
	s = "sigenu"
	return
}

// address: postgres://10.2.24.117/sigenu
// user: lamg
// pass: hqmnv78
