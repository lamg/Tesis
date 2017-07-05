package main

import (
	"github.com/lamg/tesis"
	"log"
)

func main() {
	var e error
	var s tesis.DBManager
	if e == nil {
		s, e = initDB()
	}
	e = initWindow(s, e)
	if e != nil {
		log.Fatal(e.Error())
	}
}
