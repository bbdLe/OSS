package utils

import (
	"net/http"
	"strconv"
)

func GetHashFromHeader(header http.Header) string {
	digest := header.Get("digest")
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}

	return digest[8:]
}

func GetSizeFromHeader(header http.Header) int64 {
	size, _ := strconv.ParseInt(header.Get("content-length"), 0, 64)
	return size
}

func GetOffsetFromHeader(header http.Header) int64 {
	byteRange := header.Get("range")
	size := len(byteRange)
	if size < 7 {
		return 0
	}

	if byteRange[:6] != "bytes=" {
		return 0
	}

	n, err := strconv.ParseInt(byteRange[6:size-1], 0, 64)
	if err != nil {
		return 0
	} else {
		return n
	}
}
