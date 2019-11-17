package rs

import (
	"OSS/comm/httpstream"
	"OSS/comm/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type resumableToken struct {
	Name string
	Size int64
	Hash string
	Servers []string
	Uuids []string
}

type RSResumeablePutStream struct {
	*RSPutStream
	*resumableToken
}

func NewRSResumablePutStream(dataServers []string, name string, hash string, size int64) (*RSResumeablePutStream, error) {
	putStream, err := NewRSPutStream(dataServers, hash, size)
	if err != nil {
		return nil, err
	}

	uuids := make([]string, AllShares)
	for i := range uuids {
		uuids[i] = putStream.writers[i].(*httpstream.TempPutStream).Uuid
	}

	token := &resumableToken{name, size, hash, dataServers, uuids}

	return &RSResumeablePutStream{putStream, token}, nil
}

func NewResumablePutStreamFromToken(token string) (*RSResumeablePutStream, error) {
	b, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}

	var t resumableToken
	err = json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}

	writers := make([]io.Writer, AllShares)
	for i := range writers {
		writers[i] = &httpstream.TempPutStream{t.Servers[i], t.Uuids[i]}
	}
	enc := NewEncoder(writers)
	return &RSResumeablePutStream{&RSPutStream{enc}, &t}, nil
}

func (s *RSResumeablePutStream) ToToken() string {
	b, _ := json.Marshal(s.resumableToken)
	return base64.StdEncoding.EncodeToString(b)
}

func (s *RSResumeablePutStream) CurrentSize() int64 {
	rsp, err := http.Head(fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0]))
	if err != nil {
		log.Println(err)
		return -1
	}

	if rsp.StatusCode != http.StatusOK {
		log.Println(rsp.StatusCode)
		return -1
	}

	size := utils.GetSizeFromHeader(rsp.Header) * DataShards
	if size > s.Size {
		size = s.Size
	}

	return size
}