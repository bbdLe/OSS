package httpstream

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func puthandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	b, _ := ioutil.ReadAll(r.Body)
	if string(b) != "hello" {
		w.WriteHeader(http.StatusForbidden)
	}
}

func TestPut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(puthandler))
	s, err := NewPutStream(server.URL[7:], "object")
	if err != nil {
		t.Error(err)
	}
	io.WriteString(s, "hello")
	err = s.Close()
	if err != nil {
		t.Error(err)
	}
}
