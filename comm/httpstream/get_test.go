package httpstream

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, world"))
}

func TestGet(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(handler))
	defer s.Close()

	gs, _ := NewGetStream(s.URL[7:], "objects")
	b, _ := ioutil.ReadAll(gs)
	if string(b) != "hello, world" {
		t.Errorf("except : %s, actual : %s", "hello, world", string(b))
	}
	_ = gs.Close()
}
