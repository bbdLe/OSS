package objects

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		put(w, r)
		return
	} else if r.Method == "GET" {
		get(w, r)
		return
	} else if r.Method == "DELETE" {
		del(w, r)
		return
	}


	w.WriteHeader(http.StatusInternalServerError)
}
