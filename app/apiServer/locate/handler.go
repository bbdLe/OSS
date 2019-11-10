package locate

import (
	"OSS/comm/rs"
	"encoding/json"
	"net/http"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	if len(info) < rs.DataShards {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		b, _ := json.Marshal(info)
		w.Write(b)
	}
}
