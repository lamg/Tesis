package db

import (
	"github.com/lamg/tesis"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var up, pp, ap string

func init() {
	up, pp, ap = os.Getenv("SG_USER"), os.Getenv("SG_PASS"),
		os.Getenv("SG_ADDR")
}

func TestGetStudents(t *testing.T) {
	var r tesis.RecordProvider
	var e error
	r, e = NewPSProvider(up, pp, ap, -1)
	var rs []tesis.DBRecord
	if assert.NoError(t, e) {
		rs, e = r.Records()
	}
	assert.True(t, e == nil && len(rs) > 0)
	t.Log(len(rs))
}
