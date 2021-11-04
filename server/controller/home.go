package controller

import (
	"encoding/json"
	"net/http"
)

type Home struct{}

// create a new home controller
func NewHome() *Home {
	return new(Home)
}

// get home page greeting message
func (h *Home) HomePage(rw http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	body, err := json.Marshal(map[string]interface{}{
		"message": "welcome",
	})
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(body)
}
