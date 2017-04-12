package db

import (
	"database/sql"
	a "github.com/stretchr/testify/assert"
	"testing"
)

func TestConn(t *testing.T) {
	var r *sql.Rows
	var e error
	r, e = AllStudents()
	if a.NoError(t, e) {
		r.Close()
	}
}
