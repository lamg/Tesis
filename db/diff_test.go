package db

import (
	"github.com/lamg/tesis"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestPDiff(t *testing.T) {
	var ld, sg tesis.RecordProvider
	var u, p string
	var e error
	u, p = os.Getenv("UPR_USER"), os.Getenv("UPR_PASS")
	ld, e = NewLDAPProv(u, p, "10.2.24.35:636", 300)
	if assert.NoError(t, e) {
		sg, e = NewPSProvider("lamg", "hqmnv78", "10.2.24.117/sigenu", 300)
	}
	var ds []tesis.Diff
	if assert.NoError(t, e) {
		var rp *tesis.TRpr
		rp = tesis.NewTRpr(t)
		rp.Log = false
		ds, e = PDiff(sg, ld, rp)
	}
	if assert.NoError(t, e) {
		t.Log(len(ds))
	}
}
