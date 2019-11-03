package httpstream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct{
	writer *io.PipeWriter
	ch chan error
}

func (s *PutStream) Write(p []byte) (int, error) {
	return s.writer.Write(p)
}

func (s *PutStream) Close() error {
	s.writer.Close()
	return <- s.ch
}

func newPutStream(url string) *PutStream {
	reader, writer := io.Pipe()
	ch := make(chan error)
	go func() {
		request, err := http.NewRequest("PUT", url, reader)
		if err != nil {
			reader.Close()
			ch <- err
		}
		client := http.Client{}
		r, err := client.Do(request)
		if err == nil && r.StatusCode != http.StatusOK {
			err = fmt.Errorf("DataServer return httpcode : %d", r.StatusCode)
		}
		ch <- err
	}()

	return &PutStream{writer : writer, ch : ch}
}

func NewPutStream(server string, object string) (*PutStream, error) {
	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, object)
	}

	return newPutStream(fmt.Sprintf("http://%s/objects/%s", server, object)), nil
}
