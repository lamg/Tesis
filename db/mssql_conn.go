package db

import (
	"errors"
	"github.com/lamg/tesis"
)

func NewMSSQLProvider(u, p, adr string,
	limit int) (r tesis.RecordProvider, e error) {
	e = errors.New("MSSQLProvider")
	return
}
