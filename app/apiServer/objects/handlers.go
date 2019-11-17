package objects

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPut {
		put(w, r)
		return
	} else if r.Method == http.MethodGet {
		get(w, r)
		return
	} else if r.Method == http.MethodDelete {
		del(w, r)
		return
	} else if r.Method == http.MethodPost {
		post(w, r)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
}
