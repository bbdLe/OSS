package rs

import (
	"OSS/comm/httpstream"
	"io"
)

type RSResuableGetStream struct {
	*decoder
}

func NewRSResuableGetStream(dataServers []string, uuids []string, size int64) (*RSResuableGetStream, error) {
	readers := make([]io.Reader, AllShares)
	var err error
	for i := 0; i < AllShares; i++ {
		readers[i], err = httpstream.NewTempGetStream(dataServers[i], uuids[i])
		if err != nil {
			return nil, err
		}
	}
	writers := make([]io.Writer, AllShares)
	dec := NewDecoder(readers, writers, size)
	return &RSResuableGetStream{dec}, nil
}