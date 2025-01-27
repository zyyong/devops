package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Request(method string, url string, data map[string]interface{}, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	bytesData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewReader(bytesData))
	if err != nil {
		return nil, err
	}
	for key, val := range headers {
		req.Header.Set(key, val)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}
