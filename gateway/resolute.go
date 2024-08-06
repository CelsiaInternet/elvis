package gateway

import (
	"io"
	"net/http"
	"net/url"

	"github.com/cgalvisleon/elvis/et"
	"github.com/cgalvisleon/elvis/strs"
)

type Resolute struct {
	Server     *HttpServer
	Method     string
	Proto      string
	Path       string
	RawQuery   string
	Query      url.Values
	RequestURI string
	RemoteAddr string
	Header     http.Header
	Body       io.ReadCloser
	Host       string
	Scheme     string
	Resolve    *Resolve
	URL        string
}

func GetResolute(r *http.Request) *Resolute {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	url := ""
	resolve := conn.http.GetResolve(r.Method, r.URL.Path)
	if resolve != nil {
		url = strs.Append(resolve.Resolve, r.URL.RawQuery, "?")
	}

	return &Resolute{
		Server:     conn.http,
		Method:     r.Method,
		Proto:      r.Proto,
		Path:       r.URL.Path,
		RawQuery:   r.URL.RawQuery,
		Query:      r.URL.Query(),
		RequestURI: r.RequestURI,
		RemoteAddr: r.RemoteAddr,
		Header:     r.Header,
		Body:       r.Body,
		Host:       r.Host,
		Scheme:     scheme,
		Resolve:    resolve,
		URL:        url,
	}
}

func (rs *Resolute) ToString() string {
	j := et.Json{
		"Method":     rs.Method,
		"Proto":      rs.Proto,
		"Path":       rs.Path,
		"RawQuery":   rs.RawQuery,
		"Query":      rs.Query,
		"RequestURI": rs.RequestURI,
		"RemoteAddr": rs.RemoteAddr,
		"Header":     rs.Header,
		"Body":       rs.Body,
		"Host":       rs.Host,
		"Scheme":     rs.Scheme,
		"Resolve":    rs.Resolve,
		"URL":        rs.URL,
	}

	return j.ToString()
}
