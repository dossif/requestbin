package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Request Status Handler
func (s *handler) handlerRequestStatus(w http.ResponseWriter, r *http.Request) {
	log := s.Service.Log
	httpStatus := r.PathValue("status")
	if httpStatus == "" {
		httpStatus = "200"
	}
	httpStatusInt, err := strconv.Atoi(httpStatus)
	if err != nil {
		errMsg := fmt.Sprintf("http status %v not int", httpStatus)
		log.Error(errMsg)
		jsonError(w, errMsg, 599)
		return
	}
	if httpStatusInt < 200 || httpStatusInt >= 600 {
		errMsg := fmt.Sprintf("invalid http status %v", httpStatus)
		log.Error(errMsg)
		jsonError(w, errMsg, 599)
		return
	}
	w.WriteHeader(httpStatusInt)
	body, _ := io.ReadAll(r.Body)
	cookies := make(map[string]string)
	for _, c := range r.Cookies() {
		cookies[c.Name] = c.Value
	}
	type Request struct {
		RemoteAddr    string
		Method        string
		Host          string
		Proto         string
		ProtoMajor    int
		ProtoMinor    int
		Pattern       string
		Url           string
		Path          string
		RawQuery      string
		Query         url.Values
		RequestURI    string
		ContentLength int64
		Trailer       http.Header
		Body          string
		Headers       http.Header
		Cookies       map[string]string
	}
	type Response struct {
		Status  int
		Headers http.Header
	}
	type RequestStatus struct {
		Request  Request
		Response Response
	}
	respJson := RequestStatus{
		Request: Request{
			Method:        r.Method,
			RemoteAddr:    r.RemoteAddr,
			Host:          r.Host,
			Proto:         r.Proto,
			ProtoMajor:    r.ProtoMajor,
			ProtoMinor:    r.ProtoMinor,
			Pattern:       r.Pattern,
			Url:           r.URL.String(),
			Path:          r.URL.Path,
			RawQuery:      r.URL.RawQuery,
			Query:         r.URL.Query(),
			RequestURI:    r.RequestURI,
			Trailer:       r.Trailer,
			Body:          string(body),
			Headers:       r.Header,
			Cookies:       cookies,
			ContentLength: r.ContentLength,
		},
		Response: Response{
			Status:  httpStatusInt,
			Headers: w.Header(),
		},
	}
	resp, _ := json.MarshalIndent(respJson, "", "  ")
	_, _ = w.Write(resp)
}

// Json-formatted error handler
func jsonError(w http.ResponseWriter, err string, code int) {
	w.WriteHeader(code)
	resp := struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{
		Code:  code,
		Error: err,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// Favicon handler
func (s *handler) handlerFavicon(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	_, _ = w.Write(favicon)
}
