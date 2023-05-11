package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func HttpRequest(url string, data interface{}, apikey string, requestType string) string {
	client := &http.Client{}
	bytesData, _ := json.Marshal(data)
	req, _ := http.NewRequest(requestType, url, bytes.NewReader(bytesData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", apikey)
	resp, _ := client.Do(req)
	body, _ := ioutil.ReadAll(resp.Body)
	RespData := make(map[string]string)
	json.Unmarshal(body, &RespData)
	return RespData["data"]
}
