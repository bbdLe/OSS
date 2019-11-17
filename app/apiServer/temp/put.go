package temp

import (
	"OSS/app/apiServer/config"
	"OSS/app/apiServer/locate"
	"OSS/comm/es"
	"OSS/comm/rs"
	"OSS/comm/utils"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, err := rs.NewResumablePutStreamFromToken(token)
	if err != nil {
		log.Println("NewResumablePutStreamFromToken failed :", err)
		w.WriteHeader(http.StatusForbidden)
		return
	}
	current := stream.CurrentSize()
	if current == - 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	offset := utils.GetOffsetFromHeader(r.Header)
	if offset != current {
		log.Printf("offset(%d) != current(%d)", offset, current)
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}

	bytes := make([]byte, rs.BlockSize)
	for {
		n, err := io.ReadFull(r.Body, bytes)
		if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Println("Read from stream size :", n)
		current += int64(n)
		if current > stream.Size {
			stream.Commit(false)
			log.Println("resuable put exceed size")
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if n != rs.BlockSize && current != stream.Size {
			return
		}
		stream.Write(bytes[:n])
		if current == stream.Size {
			stream.Flush()
			getStream, err := rs.NewRSResuableGetStream(stream.Servers, stream.Uuids, stream.Size)
			if err != nil {
				log.Println("GetTempStream failed :", err)
				return
			}
			h := sha256.New()
			io.Copy(h, getStream)
			calHash := base64.StdEncoding.EncodeToString(h.Sum(nil))
			if calHash != stream.Hash {
				stream.Commit(false)
				log.Println("resumable put done but hash mismatch")
				w.WriteHeader(http.StatusForbidden)
				return
			}
			if locate.Exist(stream.Hash) {
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			err = es.AddVersion(config.ServerCfg.Server.ES, stream.Name, stream.Hash, stream.Size)
			if err != nil {
				log.Println("Add Version failed :", err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}
	}
}