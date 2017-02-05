package tesis

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http"
	"io/ioutil"
)

func TestServe (t *testing.T) {
	//TODO
	l, e := serve()
	//start server
	r, e := http.Get("http://localhost:8080")
	//client make request
	assert.NoError(t, e)
	bd, e := ioutil.ReadAll(r.Body)
	assert.NoError(t, e)
	s := string(bd)
	t.Logf("s: %s", s)
	//analyze response
	//close server
	e = l.Close()
	assert.NoError(t, e)
}
