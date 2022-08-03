package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Request Status Handler
func (s *handler) handlerRequestStatus(w http.ResponseWriter, r *http.Request) {
	log := s.Service.Log
	httpStatus, _ := mux.Vars(r)["status"]
	if httpStatus == "" {
		httpStatus = "200"
	}
	httpStatusInt, err := strconv.Atoi(httpStatus)
	if err != nil {
		errMsg := fmt.Sprintf("http status %v not int", httpStatus)
		log.Errorf(errMsg)
		jsonError(w, errMsg, 599)
		return
	}
	if httpStatusInt < 200 || httpStatusInt >= 600 {
		errMsg := fmt.Sprintf("invalid http status %v", httpStatus)
		log.Errorf(errMsg)
		jsonError(w, errMsg, 599)
		return
	}
	w.WriteHeader(httpStatusInt)
	body, _ := ioutil.ReadAll(r.Body)
	type Request struct {
		RemoteAddr    string
		Method        string
		Host          string
		Proto         string
		Url           string
		RequestURI    string
		ContentLength int64
		Trailer       http.Header
		Body          string
		Headers       http.Header
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
			Url:           r.URL.String(),
			RequestURI:    r.RequestURI,
			Trailer:       r.Trailer,
			Body:          string(body),
			Headers:       r.Header,
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
	fv, _ := base64.StdEncoding.DecodeString(favicon)
	_, _ = w.Write(fv)
}
