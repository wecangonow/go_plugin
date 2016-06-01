package lib

import (
	"net/http"
	_ "net/http/pprof"
)

var Httpserver = httpserver{}

type httpserver struct {
}

func (h httpserver) StartHttp() {
	http.ListenAndServe(":8081", nil)
}

