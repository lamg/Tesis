package tesis

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
)

var password, user string

func init() {
	user = os.Getenv("UPR_USER")
	password = os.Getenv("UPR_PASS")
}

func TestConn(t *testing.T) {
	var (
		fp = "x" //arbitrary password
		fu = "y" //non registered user
	)
	_, e := conn(user, password)
	assert.NoError(t, e)

	_, e = conn(fu, fp)
	assert.Error(t, e)
}

func TestSearch(t *testing.T) {
	c, e := conn(user, password)
	if assert.NoError(t,e) {
		n, e := search(user, c)
		b := assert.NoError(t, e)
		for i := 0; b && i != len(n); i++ {
			t.Logf("Name: %s, V:%v",n[i].Name,n[i].Values)
		}
		c.Close()
	}
}
