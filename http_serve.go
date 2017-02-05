package tesis

import (
	h "net/http"
	"io"
	"net"
	"github.com/gorilla/mux"
)

/*func serveTLS() (l io.Closer, e error) {
	var (
		hp = ":10443"
		ce = "cert.pem"
		ke = "key.pem"
		dr = "/"
	)
	//TODO create Listener
	//this permits stopping the server
	//when calling to Listener.Close()
	h.HandleFunc(dr, hRoot)
	e = h.ListenAndServeTLS(hp, ce, ke, hRoot)
	return
}*/

func serve() (l io.Closer, e error) {
	var (
		nt = "tcp"
		adr = ":8080"
	)
	if ls, e := net.Listen(nt, adr); e == nil {
		//create Handler
		hr := mux.NewRouter()
		hr.HandleFunc("/", rootH).
			Methods("GET","POST")
		//register routes
		s := &h.Server {Addr:adr, Handler: hr}
		l = ls
		go func () { e = s.Serve(ls)}()
	}
	return
}

func rootH(w h.ResponseWriter, r *h.Request) {
	var (
		ct = "Content-Type"
		tp = "text/plain"
		cs = "charset"
		ut = "utf-8"
		ms = []byte("¡Hola Mundo!")
	)
	
	w.Header().Set(ct, tp)
	w.Header().Set(cs, ut)
	w.Write(ms)
}
