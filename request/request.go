package request

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/celsiainternet/elvis/et"
	"github.com/celsiainternet/elvis/utility"
)

var (
	methods = map[string]bool{
		"GET":     true,
		"POST":    true,
		"PUT":     true,
		"DELETE":  true,
		"PATCH":   true,
		"OPTIONS": true,
	}

	defaultTransport = &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Minute,
	}

	defaultClient = &http.Client{
		Timeout:   120 * time.Minute,
		Transport: defaultTransport,
	}

	bufPool = sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}
)

type Status struct {
	Ok      bool   `json:"ok"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ToJson returns a Json object
func (s Status) ToJson() et.Json {
	return et.Json{
		"ok":      s.Ok,
		"code":    s.Code,
		"message": s.Message,
	}
}

// ToString returns a string
func (s Status) ToString() string {
	return s.ToJson().ToString()
}

/**
* Body struct to convert the body response
**/
type Body struct {
	Data []byte
}

func newBody(data []byte) *Body {
	return &Body{Data: data}
}

/**
* ToJson returns a Json object
* @return et.Json
**/
func (b Body) ToJson() (et.Json, error) {
	var result et.Json
	err := json.Unmarshal(b.Data, &result)
	if err != nil {
		return et.Json{}, err
	}

	return result, nil
}

/**
* ToItem returns a Item object
* @return et.Item
**/
func (b Body) ToItem() (et.Item, error) {
	var result et.Item
	err := json.Unmarshal(b.Data, &result)
	if err != nil {
		return et.Item{}, err
	}

	return result, nil
}

/**
* ToItems returns a Items object
* @return et.Items
**/
func (b Body) ToItems() (et.Items, error) {
	var result et.Items
	err := json.Unmarshal(b.Data, &result)
	if err != nil {
		return et.Items{}, err
	}

	return result, nil
}

/**
* ToArrayJson returns a Json array object
* @return []et.Json
**/
func (b Body) ToArrayJson() ([]et.Json, error) {
	var result []et.Json
	err := json.Unmarshal(b.Data, &result)
	if err != nil {
		return []et.Json{}, err
	}

	return result, nil
}

/**
* ToString returns a string
* @return string
**/
func (b Body) ToString() string {
	return string(b.Data)
}

/**
* ToInt returns an integer
* @return int
**/
func (b Body) ToInt() (int, error) {
	var result int
	err := json.Unmarshal(b.Data, &result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

/**
* ToInt64 returns an integer
* @return int64
**/
func (b Body) ToInt64() (int64, error) {
	var result int64
	err := json.Unmarshal(b.Data, &result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

/**
* ToFloat returns a float
* @return float64
**/
func (b Body) ToFloat() (float64, error) {
	var result float64
	err := json.Unmarshal(b.Data, &result)
	if err != nil {
		return 0, err
	}

	return result, nil
}

/**
* ToBool returns a boolean
* @return bool
**/
func (b Body) ToBool() (bool, error) {
	var result bool
	err := json.Unmarshal(b.Data, &result)
	if err != nil {
		return false, err
	}

	return result, nil
}

/**
* ToTime returns a time
* @return time.Time
**/
func (b Body) ToTime() (time.Time, error) {
	var result time.Time
	err := json.Unmarshal(b.Data, &result)
	if err != nil {
		return time.Time{}, err
	}

	return result, nil
}

/**
* ReadBody reads the body response
* @param body io.ReadCloser
* @return *Body, error
**/
func ReadBody(body io.ReadCloser) (*Body, error) {
	bodyBytes, err := io.ReadAll(body)
	if err != nil {
		return newBody([]byte("")), err
	}

	return newBody(bodyBytes), nil
}

/**
* statusOk
* @param status int
* @return bool
**/
func statusOk(status int) bool {
	return status >= http.StatusOK && status < http.StatusMultipleChoices
}

/**
* bodyParams
* @param header, body et.Json
* @return []byte
**/
func bodyParams(header, body et.Json) []byte {
	contentType := header.Get("Content-Type")
	if contentType == "application/x-www-form-urlencoded" {
		data := url.Values{}
		for k, v := range body {
			data.Set(k, v.(string))
		}
		return []byte(data.Encode())
	} else if contentType == "application/json" {
		return []byte(body.ToEscapeHTML())
	} else {
		return []byte(body.ToString())
	}
}

type HttpResult struct {
	Body   *Body
	Status Status
}

/**
* httpDo executes an HTTP request with context and buffer pool support.
* @param ctx context.Context, method, path string, header, body et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func httpDo(ctx context.Context, method, path string, header, body et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	result := make(chan HttpResult, 1)

	go func() {
		if _, ok := methods[method]; !ok {
			result <- HttpResult{
				Body: nil,
				Status: Status{
					Ok:      false,
					Code:    http.StatusBadRequest,
					Message: "Invalid method",
				},
			}
			return
		}

		contentType := header.Str("Content-Type")

		mediaType, _, err := mime.ParseMediaType(contentType)
		if err != nil {
			result <- HttpResult{
				Body: nil,
				Status: Status{
					Ok:      false,
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			}
			return
		}

		var ioBody io.Reader
		var buf *bytes.Buffer

		switch mediaType {
		case "multipart/form-data":
			writer := multipart.NewWriter(buf)
			for k := range body {
				v := body.Str(k)
				if err := writer.WriteField(k, v); err != nil {
					result <- HttpResult{
						Body: nil,
						Status: Status{
							Ok:      false,
							Code:    http.StatusBadRequest,
							Message: err.Error(),
						},
					}
					return
				}
			}
			if err := writer.Close(); err != nil {
				result <- HttpResult{
					Body: nil,
					Status: Status{
						Ok:      false,
						Code:    http.StatusBadRequest,
						Message: err.Error(),
					},
				}
				return
			}
			ioBody = buf
		case "application/x-www-form-urlencoded":
			data := url.Values{}
			for k := range body {
				v := body.Str(k)
				data.Set(k, v)
			}
			ioBody = bytes.NewBufferString(data.Encode())
		case "application/json":
			if body != nil {
				buf = bufPool.Get().(*bytes.Buffer)
				buf.Reset()
				buf.Write(bodyParams(header, body))
				ioBody = buf
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, path, ioBody)
		if err != nil {
			if buf != nil {
				bufPool.Put(buf)
			}

			result <- HttpResult{
				Body: nil,
				Status: Status{
					Ok:      false,
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			}
			return
		}

		for k, v := range header {
			req.Header.Set(k, v.(string))
		}

		client := defaultClient
		if tlsConfig != nil {
			client = &http.Client{
				Timeout: defaultClient.Timeout,
				Transport: &http.Transport{
					TLSClientConfig:     tlsConfig,
					MaxIdleConns:        defaultTransport.MaxIdleConns,
					MaxIdleConnsPerHost: defaultTransport.MaxIdleConnsPerHost,
					IdleConnTimeout:     defaultTransport.IdleConnTimeout,
				},
			}
		}

		res, err := client.Do(req)
		if buf != nil {
			bufPool.Put(buf)
		}

		if err != nil {
			result <- HttpResult{
				Body: nil,
				Status: Status{
					Ok:      false,
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			}
			return
		}
		defer res.Body.Close()

		resultBody, err := ReadBody(res.Body)
		if err != nil {
			result <- HttpResult{
				Body: nil,
				Status: Status{
					Ok:      false,
					Code:    http.StatusBadRequest,
					Message: err.Error(),
				},
			}
			return
		}

		result <- HttpResult{
			Body: resultBody,
			Status: Status{
				Ok:      statusOk(res.StatusCode),
				Code:    res.StatusCode,
				Message: res.Status,
			},
		}
	}()

	if timeout == 0 {
		select {
		case respuesta := <-result:
			return respuesta.Body, respuesta.Status
		}
	} else {
		select {
		case respuesta := <-result:
			// La función terminó antes de timeout
			return respuesta.Body, respuesta.Status

		case <-time.After(timeout):
			// Se agotaron los timeout
			return newBody(defaultValue), Status{
				Ok:      false,
				Code:    http.StatusRequestTimeout,
				Message: "timeout",
			}
		}
	}
}

/**
* HttpCtx executes an HTTP request honoring the provided context for cancellation and deadlines.
* @param ctx context.Context, method, path string, header, body et.Json
* @return *Body, Status
**/
func HttpCtx(ctx context.Context, method, path string, header, body et.Json, tlsConfig *tls.Config) (*Body, Status) {
	return httpDo(ctx, method, path, header, body, tlsConfig, 0, nil)
}

/**
* Http
* @param method, path string, header, body et.Json
* @return *Body, Status
**/
func Http(method, path string, header, body et.Json, tlsConfig *tls.Config) (*Body, Status) {
	return httpDo(context.Background(), method, path, header, body, tlsConfig, 0, nil)
}

/**
* Post
* @param path string, header, body et.Json
* @return *Body, Status
**/
func Post(path string, header, body et.Json) (*Body, Status) {
	return Http("POST", path, header, body, nil)
}

/**
* Get
* @param path string, header et.Json
* @return *Body, Status
**/
func Get(path string, header et.Json) (*Body, Status) {
	return Http("GET", path, header, nil, nil)
}

/**
* Put
* @param path string, header, body et.Json
* @return *Body, Status
**/
func Put(path string, header, body et.Json) (*Body, Status) {
	return Http("PUT", path, header, body, nil)
}

/**
* Delete
* @param path string, header et.Json
* @return *Body, Status
**/
func Delete(path string, header et.Json) (*Body, Status) {
	return Http("DELETE", path, header, et.Json{}, nil)
}

/**
* Patch
* @param path string, header, body et.Json
* @return *Body, Status
**/
func Patch(path string, header, body et.Json) (*Body, Status) {
	return Http("PATCH", path, header, body, nil)
}

/**
* Options
* @param path string, header et.Json
* @return *Body, Status
**/
func Options(path string, header et.Json) (*Body, Status) {
	return Http("OPTIONS", path, header, et.Json{}, nil)
}

/**
* PostWithTls
* @param path string, header, body et.Json, tlsConfig *tls.Config
* @return *Body, Status
**/
func PostWithTls(path string, header, body et.Json, tlsConfig *tls.Config) (*Body, Status) {
	return Http("POST", path, header, body, tlsConfig)
}

/**
* GetWithTls
* @param path string, header et.Json, tlsConfig *tls.Config
* @return *Body, Status
**/
func GetWithTls(path string, header et.Json, tlsConfig *tls.Config) (*Body, Status) {
	return Http("GET", path, header, et.Json{}, tlsConfig)
}

/**
* PutWithTls
* @param path string, header, body et.Json, tlsConfig *tls.Config
* @return *Body, Status
**/
func PutWithTls(path string, header, body et.Json, tlsConfig *tls.Config) (*Body, Status) {
	return Http("PUT", path, header, body, tlsConfig)
}

/**
* DeleteWithTls
* @param path string, header et.Json, tlsConfig *tls.Config
* @return *Body, Status
**/
func DeleteWithTls(path string, header et.Json, tlsConfig *tls.Config) (*Body, Status) {
	return Http("DELETE", path, header, et.Json{}, tlsConfig)
}

/**
* PatchWithCA
* @param path string, header, body et.Json, tlsConfig *tls.Config
* @return *Body, Status
**/
func PatchWithTls(path string, header, body et.Json, tlsConfig *tls.Config) (*Body, Status) {
	return Http("PATCH", path, header, body, tlsConfig)
}

/**
* OptionsWithTls
* @param path string, header et.Json, tlsConfig *tls.Config
* @return *Body, Status
**/
func OptionsWithTls(path string, header et.Json, tlsConfig *tls.Config) (*Body, Status) {
	return Http("OPTIONS", path, header, et.Json{}, tlsConfig)
}

/**
* HttpCtxWithTimeout executes an HTTP request honoring the provided context for cancellation and deadlines.
* @param ctx context.Context, method, path string, header, body et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func HttpCtxWithTimeout(ctx context.Context, method, path string, header, body et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(ctx, method, path, header, body, tlsConfig, timeout, defaultValue)
}

/**
* HttpWithTimeout executes an HTTP request honoring the provided timeout.
* @param method, path string, header, body et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func HttpWithTimeout(method, path string, header, body et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), method, path, header, body, tlsConfig, timeout, defaultValue)
}

/**
* PostWithTlsTimeout executes a POST request honoring the provided timeout.
* @param path string, header, body et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func PostWithTlsTimeout(path string, header, body et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), "POST", path, header, body, tlsConfig, timeout, defaultValue)
}

/**
* GetWithTlsTimeout executes a GET request honoring the provided timeout.
* @param path string, header et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func GetWithTlsTimeout(path string, header et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), "GET", path, header, et.Json{}, tlsConfig, timeout, defaultValue)
}

/**
* PutWithTlsTimeout executes a PUT request honoring the provided timeout.
* @param path string, header, body et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func PutWithTlsTimeout(path string, header, body et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), "PUT", path, header, body, tlsConfig, timeout, defaultValue)
}

/**
* DeleteWithTlsTimeout executes a DELETE request honoring the provided timeout.
* @param path string, header et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func DeleteWithTlsTimeout(path string, header et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), "DELETE", path, header, et.Json{}, tlsConfig, timeout, defaultValue)
}

/**
* PatchWithTlsTimeout executes a PATCH request honoring the provided timeout.
* @param path string, header, body et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func PatchWithTlsTimeout(path string, header, body et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), "PATCH", path, header, body, tlsConfig, timeout, defaultValue)
}

/**
* OptionsWithTlsTimeout executes a OPTIONS request honoring the provided timeout.
* @param path string, header et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func OptionsWithTlsTimeout(path string, header et.Json, tlsConfig *tls.Config, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), "OPTIONS", path, header, et.Json{}, tlsConfig, timeout, defaultValue)
}

/**
* PostWithTimeout executes a POST request honoring the provided timeout.
* @param path string, header, body et.Json, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func PostWithTimeout(path string, header, body et.Json, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), "POST", path, header, body, nil, timeout, defaultValue)
}

/**
* GetWithTimeout executes a GET request honoring the provided timeout.
* @param path string, header et.Json, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func GetWithTimeout(path string, header et.Json, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), "GET", path, header, et.Json{}, nil, timeout, defaultValue)
}

/**
* PutWithTimeout executes a PUT request honoring the provided timeout.
* @param path string, header, body et.Json, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func PutWithTimeout(path string, header, body et.Json, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), "PUT", path, header, body, nil, timeout, defaultValue)
}

/**
* DeleteWithTimeout executes a DELETE request honoring the provided timeout.
* @param path string, header et.Json, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func DeleteWithTimeout(path string, header et.Json, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), "DELETE", path, header, et.Json{}, nil, timeout, defaultValue)
}

/**
* PatchWithTimeout executes a PATCH request honoring the provided timeout.
* @param path string, header, body et.Json, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func PatchWithTimeout(path string, header, body et.Json, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), "PATCH", path, header, body, nil, timeout, defaultValue)
}

/**
* OptionsWithTimeout executes a OPTIONS request honoring the provided timeout.
* @param path string, header et.Json, timeout time.Duration, defaultValue []byte
* @return *Body, Status
**/
func OptionsWithTimeout(path string, header et.Json, timeout time.Duration, defaultValue []byte) (*Body, Status) {
	return httpDo(context.Background(), "OPTIONS", path, header, et.Json{}, nil, timeout, defaultValue)
}

/**
* NewTlsConfig
* @param caPath, certPath, keyPath string
* @return *tls.Config, error
**/
func NewTlsConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
	if certFile == "" {
		return nil, fmt.Errorf("CRT certificate path is required")
	}

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("CRT certificate not found")
	}

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	if !utility.ValidStr(caFile, 0, []string{""}) {
		return tlsConfig, nil
	}

	caCert, err := os.ReadFile(caFile)
	if !os.IsNotExist(err) {
		tlsConfig.RootCAs = x509.NewCertPool()
		tlsConfig.RootCAs.AppendCertsFromPEM(caCert)
	}

	return tlsConfig, nil
}
