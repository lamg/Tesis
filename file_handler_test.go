package tesis

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFileHandler(t *testing.T) {
	var fl *FileHandler
	var e error
	var nm string
	nm = "ok.file"
	fl, e = NewFileHandler(nm)
	assert.NoError(t, e)
	_, e = fl.Write([]byte("¡Hola Mundo!"))
	if assert.NoError(t, e) {
		fl.Close()
		e = os.Remove(nm)
	}
	assert.NoError(t, e)
}
