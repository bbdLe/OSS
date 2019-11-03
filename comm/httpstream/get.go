package httpstream

import (
	"fmt"
	"io"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

func newGetStream(url string) (*GetStream, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer return http code %d", resp.StatusCode)
	}
	return &GetStream{reader : resp.Body}, nil
}

func NewGetStream(server string, object string) (*GetStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("invaild server %s object %s", server, object)
	}
	return newGetStream(fmt.Sprintf("http://%s/objects/%s", server, object))
}

func (s *GetStream) Read(p []byte) (int, error) {
	return s.reader.Read(p)
}
