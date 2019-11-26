package rs

import (
	"OSS/comm/httpstream"
	"fmt"
	"io"
)

type RSPutStream struct {
	*encoder
}

func NewRSPutStream(dataServer []string, hash string, size int64) (*RSPutStream, error) {
	if len(dataServer) != AllShares {
		return nil, fmt.Errorf("dataServers number mismatch")
	}

	perShard := (size + DataShards - 1) / DataShards
	writers := make([]io.Writer, AllShares)
	for i := range writers {
		var err error
		writers[i], err = httpstream.NewTempPutStream(dataServer[i], fmt.Sprintf("%s.%d", hash, i), perShard)
		if err != nil {
			return nil, err
		}
	}
	enc := NewEncoder(writers)

	return &RSPutStream{enc}, nil
}

func (s *RSPutStream) Commit(succ bool) {
	s.Flush()
	for i := range s.writers {
		s.writers[i].(*httpstream.TempPutStream).Commit(succ)
	}
	s.Close()
}
