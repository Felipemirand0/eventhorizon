package controller

import (
	"net/http"
)

type httpHealthz struct {
}

func (h httpHealthz) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("ok"))
}
