package controller

import (
	"net/http"
)

type Stopper struct {
	c chan struct{}
}

// create a new stopper controller
func NewStopper(c chan struct{}) *Stopper {
	return &Stopper{
		c: c,
	}
}

// api for stop server
func (h *Stopper) StopServer(rw http.ResponseWriter, r *http.Request) {
	// limit calling method
	if r.Method != "GET" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	h.c <- struct{}{}
}
