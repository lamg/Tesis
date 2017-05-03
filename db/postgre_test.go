package db

import (
	"github.com/lamg/tesis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetStudents(t *testing.T) {
	var r tesis.RecordProvider
	var e error
	r, e = NewPSProvider("lamg", "hqmnv78", "10.2.24.117/sigenu", -1)
	var rs []tesis.DBRecord
	if assert.NoError(t, e) {
		rs, e = r.Records()
	}
	if assert.NoError(t, e) {
		t.Log(len(rs))
	}
}
