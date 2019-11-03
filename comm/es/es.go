package es

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	url2 "net/url"
	"strings"
)

type Metadata struct {
	Name string
	Version int
	Size int64
	Hash string
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Hits []hit
	}
}

func getMetaData(server string, name string, version int) (meta Metadata, err error) {
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d/_source",
		server, name, version)
	r, err := http.Get(url)
	if err != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to get %s_%d: %d", name, version, r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(result, &meta)
	return
}

func PutMetaData(server string, name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":%d,"size":%d,"hash":"%s"}`,
		name, version, size, hash)
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d?op_type=create", server, name, version)

	request, err := http.NewRequest("POST", url, strings.NewReader(doc))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	r, e := client.Do(request)
	if e != nil {
		return e
	}
	if r.StatusCode == http.StatusConflict {
		return PutMetaData(server, name, version + 1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata : %d %s", r.StatusCode, result)
	}

	return nil
}

func SearchLatestVersion(server string, name string) (meta Metadata, err error){
	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc",
		server, url2.PathEscape(name))
	r, err := http.Get(url)
	if err != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(r.Body)
		err = fmt.Errorf("fail to search lastest metadata : (%d) %v", r.StatusCode,
			string(msg))
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	err = json.Unmarshal(result, &sr)
	if err != nil {
		return
	}
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source
	}
	return
}

func SearchAllVersion(server string, name string, from int, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?=sort=name,version&from=%d&size=%d",
		server, from, size)
	if name != "" {
		url += "&q=name:" + name
	}
	log.Println(url)
	r, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get all version : %d", r.StatusCode)
	}

	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	err = json.Unmarshal(result, &sr)
	if err != nil {
		return nil, err
	}
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)
	}

	return metas, nil
}

func GetMetadata(server string, name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(server, name)
	} else {
		return getMetaData(server, name, version)
	}
}

func DelMetadata(server string, name string, version int) {
	url := fmt.Sprintf("http://%s/metadata/objects/%s_%d",
		server, name, version)
	client := http.Client{}
	request, _ := http.NewRequest("DELETE", url, nil)
	client.Do(request)
}

func AddVersion(server string, name string, hash string, size int64) error {
	data, err := GetMetadata(server, name, 0)
	if err != nil {
		return err
	}
	err = PutMetaData(server, name, data.Version + 1, size, hash)
	return err
}