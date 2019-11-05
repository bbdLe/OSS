package temp

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		put(w, r)
		return
	}
	if r.Method == "POST" {
		post(w, r)
		return
	}
	if r.Method == "PATCH" {
		patch(w, r)
		return
	}
}