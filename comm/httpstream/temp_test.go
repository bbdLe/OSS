package httpstream

import (
	"testing"
)

func TestPutStream(t *testing.T) {
	stream, err := NewTempPutStream("localhost:5001",
		"aWKQ2BipX94sb+h3xdTbWYAu1yzjn5vyFG2SOwUQIXY=", 37)
	if err != nil {
		t.Error(err)
	}

	_, err = stream.Write([]byte("this object will have only 1 instance"))
	if err != nil {
		t.Error(err)
	}
	stream.Commit(true)
}
