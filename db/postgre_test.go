package db

import (
	"github.com/lamg/tesis"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetStudents(t *testing.T) {
	var r tesis.RecordProvider
	var e error
	r, e = NewPSProvider("lamg", "hqmnv78")
	var rs []tesis.DBRecord
	if assert.NoError(t, e) {
		rs, e = r.Records()
	}
	if assert.NoError(t, e) {
		t.Log(len(rs))
	}
}
