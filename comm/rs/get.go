package rs

import (
	"OSS/comm/httpstream"
	"fmt"
	"io"
)

type RSGetStream struct {
	*decoder
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string,
	hash string, size int64) (*RSGetStream, error) {
		if len(locateInfo) + len(dataServers) != AllShares {
			return nil, fmt.Errorf("dataServers number dismatch")
		}

		readers := make([]io.Reader, AllShares)
		for i := 0; i < AllShares; i++ {
			server := locateInfo[i]
			if server == "" {
				locateInfo[i] = dataServers[0]
				dataServers = dataServers[1:]
				continue
			}
			reader, err := httpstream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
			if err == nil {
				readers[i] = reader
			}
		}

		writers := make([]io.Writer, AllShares)
		perShard := (size + DataShards - 1) / DataShards
		var err error
		for i := 0; i < AllShares; i++ {
			if readers[i] == nil {
				writers[i], err = httpstream.NewTempPutStream(locateInfo[i],
					fmt.Sprintf("%s.%d", hash, i), perShard)
				if err != nil {
					return nil, err
				}
			}
		}

		decoder := NewDecoder(readers, writers, size)
		return &RSGetStream{decoder}, nil
}

func (s *RSGetStream) Close() {
	for i := range s.writers {
		if s.writers[i] != nil {
			s.writers[i].(*httpstream.TempPutStream).Commit(true)
		}
	}
}