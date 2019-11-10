package rs

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

func testEncodeDecode(t *testing.T, p []byte) {
	writers := make([]io.Writer, AllShares)
	readers := make([]io.Reader, AllShares)
	for i := range writers {
		writers[i], _ = os.Create(fmt.Sprintf("/tmp/ut_%d", i))
	}
	enc := NewEncoder(writers)
	length := len(p)
	for count := 0; count != length; {
		n, e := enc.Write(p[count:])
		if e != nil {
			t.Error(e)
		}
		count += n
	}
	enc.Flush()
	enc.Close()
	for i := range writers {
		writers[i] = nil
		readers[i], _ = os.Open(fmt.Sprintf("/tmp/ut_%d", i))
	}
	readers[1] = nil
	readers[4] = nil
	writers[1], _ = os.Create(fmt.Sprintf("/tmp/repair_1"))
	writers[4], _ = os.Create(fmt.Sprintf("/tmp/repair_4"))
	dec := NewDecoder(readers, writers, int64(length))
	b := make([]byte, length + 10)
	count := 0
	for {
		n, e := dec.Read(b)
		count += n
		if e == io.EOF {
			break
		}
	}
	if count != length {
		t.Error(count, length)
	}
	if !reflect.DeepEqual(b[:count], p) {
		t.Errorf("not match")
	}
	output, e := exec.Command("diff", "/tmp/ut_1", "/tmp/repair_1").Output()
	if len(output) != 0 {
		t.Error(string(output), e)
	}

	output, e = exec.Command("diff", "/tmp/ut_4", "/tmp/repair_4").Output()
	if len(output) != 0 {
		t.Error(string(output), e)
	}
}

func TestEncodeDecode(t *testing.T) {
	//p := []byte{1}
	//testEncodeDecode(t, p)
	p := []byte("1111111111111111111111111")
	testEncodeDecode(t, p)
}
