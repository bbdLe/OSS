package es

import (
	"testing"
	"time"
)

const server = "localhost:9200"

const version = 3

func TestPutMetaData(t *testing.T) {
	err := PutMetaData(server, "test", 1, 0, "")
	if err != nil {
		t.Error(err)
	}
}

func TestSearchLatestVersion(t *testing.T) {
	time.Sleep(time.Second * 5)
	meta, err := SearchLatestVersion(server, "test")
	if err != nil {
		t.Error(err)
	}
	if meta.Version < 1 {
		t.Errorf("Wront Version")
	}
}
