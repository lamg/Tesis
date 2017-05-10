package db

import (
	"database/sql"
	"fmt"
	"github.com/lamg/tesis"
	_ "github.com/lib/pq"
)

func NewPSProvider(u, p, adr string, t int) (r tesis.RecordProvider, e error) {
	var l *PSProv
	l = new(PSProv)
	var s string
	s = fmt.Sprintf("postgres://%s:%s@%s", u, p, adr)
	s = "postgres://lamg:hqmnv78@10.2.24.117/sigenu" //testing database
	l.db, e = sql.Open("postgres", s)
	if e == nil {
		l.limit = t
		r = l
	}
	return
}

type PSProv struct {
	db    *sql.DB
	limit int
}

func (p *PSProv) Records() (s []tesis.DBRecord, e error) {
	var r *sql.Rows
	var qr string
	if p.limit >= 0 {
		qr = fmt.Sprintf("SELECT id_student,identification,name,middle_name,last_name,address,phone FROM student LIMIT %d", p.limit)
	} else {
		qr = "SELECT id_student,identification,name,middle_name,last_name,address,phone FROM student where student_status_fk='02'"
	}
	r, e = p.db.Query(qr)
	s = make([]tesis.DBRecord, 0)
	for e == nil && r.Next() {
		var st tesis.DBRecord
		st = tesis.DBRecord{}
		var name, middle_name, last_name string
		e = r.Scan(&st.Id, &st.IN, &name, &middle_name, &last_name, &st.Addr, &st.Tel)
		if e == nil {
			st.Name = name + " " + middle_name + " " + last_name
			s = append(s, st)
		}
	}
	return
}

func (p *PSProv) Name() (s string) {
	s = "sigenu"
	return
}
