package httpstream

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type TempPutStream struct {
	Server string
	Uuid string
}

func NewTempPutStream(server, object string, size int64) (*TempPutStream, error) {
	request, err := http.NewRequest("POST", "http://" + server + "/temp/" + object, nil)
	if err != nil {
		log.Println("Request failed :", err)
		return nil, err
	}

	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Client do failed :", err)
		return nil, err
	}

	uuid, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("read response body failed :", err)
		return nil, err
	}

	return &TempPutStream{server, string(uuid)}, nil
}

func (w *TempPutStream) Write(p []byte) (n int, err error) {
	request, err := http.NewRequest("PATCH", "http://" + w.Server + "/temp/" + w.Uuid, strings.NewReader(string(p)))
	if err != nil {
		log.Println("Request fail :", request)
		return
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Client do failed :", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("response status : %d", response.StatusCode)
		return
	}

	return len(p), nil
}

func (w *TempPutStream) Commit(succ bool) {
	method := "DELETE"
	if succ {
		method = "PUT"
	}

	request, err := http.NewRequest(method, "http://" + w.Server + "/temp/" + w.Uuid, nil)
	if err != nil {
		log.Println("NewRequest failed ", err)
		return
	}
	client := http.Client{}
	_, err = client.Do(request)
	if err != nil {
		log.Println("Client do fail ", err)
		return
	}
}

func NewTempGetStream(server, uuid string) (*GetStream, error) {
	return newGetStream("http://" + server + "/temp/" + uuid)
}
