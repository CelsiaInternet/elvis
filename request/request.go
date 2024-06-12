package request

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/cgalvisleon/elvis/console"
	"github.com/cgalvisleon/elvis/et"
)

// ioReadeToJson reads the io.Reader and returns a Json object
func ioReadeToJson(r io.Reader) (et.Json, error) {
	var result et.Json
	err := json.NewDecoder(r).Decode(&result)
	if err != nil {
		return et.Json{}, err
	}

	return result, nil
}

// Request post method
func Post(url string, header, body et.Json) (et.Json, error) {
	bodyParams := []byte(body.ToString())
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyParams))
	if err != nil {
		return et.Json{}, err
	}

	for k, v := range header {
		req.Header.Set(k, v.(string))
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return et.Json{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return et.Json{}, console.ErrorF(`%s - Status:%d`, res.Status, res.StatusCode)
	}

	result, err := ioReadeToJson(res.Body)
	if err != nil {
		return et.Json{}, nil
	}

	return result, nil

}

// Request get method
func Get(url string, header et.Json) (et.Json, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return et.Json{}, err
	}

	for k, v := range header {
		req.Header.Set(k, v.(string))
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return et.Json{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return et.Json{}, console.ErrorF(`%s - Status:%d`, res.Status, res.StatusCode)
	}

	result, err := ioReadeToJson(res.Body)
	if err != nil {
		return et.Json{}, nil
	}

	return result, nil
}

// Request put method
func Put(url string, header, body et.Json) (et.Json, error) {
	bodyParams := []byte(body.ToString())
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(bodyParams))
	if err != nil {
		return et.Json{}, err
	}

	for k, v := range header {
		req.Header.Set(k, v.(string))
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return et.Json{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return et.Json{}, console.ErrorF(`%s - Status:%d`, res.Status, res.StatusCode)
	}

	result, err := ioReadeToJson(res.Body)
	if err != nil {
		return et.Json{}, nil
	}

	return result, nil
}

// Request delete method
func Delete(url string, header et.Json) (et.Json, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return et.Json{}, err
	}

	for k, v := range header {
		req.Header.Set(k, v.(string))
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return et.Json{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return et.Json{}, console.ErrorF(`%s - Status:%d`, res.Status, res.StatusCode)
	}

	result, err := ioReadeToJson(res.Body)
	if err != nil {
		return et.Json{}, nil
	}

	return result, nil
}

// Request patch method
func Patch(url string, header, body et.Json) (et.Json, error) {
	bodyParams := []byte(body.ToString())
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(bodyParams))
	if err != nil {
		return et.Json{}, err
	}

	for k, v := range header {
		req.Header.Set(k, v.(string))
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return et.Json{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return et.Json{}, console.ErrorF(`%s - Status:%d`, res.Status, res.StatusCode)
	}

	result, err := ioReadeToJson(res.Body)
	if err != nil {
		return et.Json{}, nil
	}

	return result, nil
}

// Request options method
func Options(url string, header et.Json) (et.Json, error) {
	req, err := http.NewRequest("OPTIONS", url, nil)
	if err != nil {
		return et.Json{}, err
	}

	for k, v := range header {
		req.Header.Set(k, v.(string))
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return et.Json{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return et.Json{}, console.ErrorF(`%s - Status:%d`, res.Status, res.StatusCode)
	}

	result, err := ioReadeToJson(res.Body)
	if err != nil {
		return et.Json{}, nil
	}

	return result, nil
}
