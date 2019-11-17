package temp

import (
	"OSS/comm/rs"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func head(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := rs.NewResumablePutStreamFromToken(token)
	if err != nil {
		log.Println("NewResumablePutStreamFromToken failed :", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == -1 {
		log.Println("failed to get size")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("content-length", fmt.Sprintf("%d", current))
}
