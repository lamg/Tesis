package db

import (
	"github.com/lamg/tesis"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSync(t *testing.T) {
	var ld, sg tesis.RecordProvider
	var u, p string
	var e error
	u, p = os.Getenv("UPR_USER"), os.Getenv("UPR_PASS")
	ld, e = NewLDAPProv(u, p)
	if assert.NoError(t, e) {
		sg, e = NewPSProvider("lamg", "hqmnv78")
	}
	var ds []tesis.Diff
	if assert.NoError(t, e) {
		ds, e = Sync(sg, ld)
	}
	if assert.NoError(t, e) {
		t.Log(len(ds))
	}
}
