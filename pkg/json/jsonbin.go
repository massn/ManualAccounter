package jsonbin

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

const route = "https://api.jsonbin.io/v3/b"

type CreateParam struct {
	BinName   string
	Body      string
	IsPrivate bool
	APIKey    string
}

type Response struct {
	StatusCode int
	Metadata   struct {
		Id        string `json:"id"`
		CreatedAt string `json:"createdAt"`
		Private   bool   `json:"private"`
		Name      string `json:"name"`
	} `json:"metadata"`
	Record  json.RawMessage
	Message string `json:"message"`
}

func Create(cp CreateParam) (*Response, error) {
	req, err := http.NewRequest("POST", route, bytes.NewBuffer([]byte(cp.Body)))
	if err != nil {
		return &Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Master-Key", cp.APIKey)
	req.Header.Set("X-Bin-Name", cp.BinName)
	req.Header.Set("X-Bin-Private", strconv.FormatBool(cp.IsPrivate))

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return &Response{}, err
	}
	defer resp.Body.Close()
	return makeResponse(resp)
}

type ReadParam struct {
	BinId      string
	BinVersion string
	APIKey     string
}

func Read(rp ReadParam) (*Response, error) {
	url := route + "/" + rp.BinId
	if rp.BinVersion != "" {
		url = url + "/" + rp.BinVersion
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return &Response{}, err
	}
	req.Header.Set("X-Master-Key", rp.APIKey)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return &Response{}, err
	}
	defer resp.Body.Close()
	return makeResponse(resp)
}

func makeResponse(resp *http.Response) (*Response, error) {
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &Response{}, err
	}

	res := Response{StatusCode: resp.StatusCode}
	if err := json.Unmarshal(byteArray, &res); err != nil {
		return &Response{}, err
	}
	return &res, nil
}
