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
		r, e = db.Query("SELECT identification,name FROM student")
	}
	ss = make([]tesis.DBRecord, 0)
	for r.Next() {
		var st tesis.DBRecord
		st = tesis.DBRecord{}
		r.Scan(&st.Id, &st.Name)
		ss = append(ss, st)
	}
	return
}

// address: postgres://10.2.24.117/sigenu
// user: lamg
// pass: hqmnv78
