package db

import (
	"database/sql"
	"github.com/lamg/tesis"
	_ "github.com/lib/pq"
)

func AllStudents() (ss []tesis.DBRecord, e error) {
	var db *sql.DB
	var r *sql.Rows
	db, e = sql.Open("postgres", "postgres://lamg:hqmnv78@10.2.24.117/sigenu")
	if e == nil {
		r, e = db.Query("SELECT id_student,identification,name,middle_name,last_name FROM student")
	}
	ss = make([]tesis.DBRecord, 0)
	for e == nil && r.Next() {
		var st tesis.DBRecord
		st = tesis.DBRecord{}
		var name, middle_name, last_name string
		e = r.Scan(&st.Id, &st.IN, &name, &middle_name, &last_name)
		if e == nil {
			st.Name = name + middle_name + last_name
			ss = append(ss, st)
		}
	}
	return
}

// address: postgres://10.2.24.117/sigenu
// user: lamg
// pass: hqmnv78
