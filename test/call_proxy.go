package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/janqii/go-httpclient"
	"io/ioutil"
	"net/http"
	"time"
)

type UserData struct {
	uid     int
	uname   string
	content string
}
type Request struct {
	Topic        string
	PartitionKey string
	Data         UserData
}

func main() {
	udata := UserData{
		uid:     123,
		uname:   "crask",
		content: "welcome to crask",
	}

	rdata := Request{
		Topic:        "test",
		PartitionKey: "123",
		Data:         udata,
	}

	url := "http://127.0.0.1:9090/produce?format=json"

	ret, err := _httpCall(url, rdata, 123456)
	if err != nil {
		fmt.Printf("call proxy error: %v\n", err)
	} else {
		fmt.Printf("call proxy ok: %v\n", ret)
	}
}

func _httpCall(url string, input interface{}, logId int) (string, error) {
	if "" == url {
		return "", errors.New("invalid url")
	}

	s, e := json.Marshal(input)
	if e != nil {
		return "", e
	}

	b := []byte(s)
	r := bytes.NewReader(b)

	transport := &httpclient.Transport{
		ConnectTimeout:        1 * time.Second,
		RequestTimeout:        1 * time.Second,
		ResponseHeaderTimeout: 1 * time.Second,
	}
	defer transport.Close()

	client := &http.Client{Transport: transport}

	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return "", err
	}

	req.Header.Add("X_BD_LOGID", fmt.Sprintf("%d", logId))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return string(body), e
	}

	return string(body), e
}
