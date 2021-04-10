package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

// Split full-address to address and port
func splitAddrPort(host string) (addr string, port int) {
	lastInd := strings.LastIndex(host, ":")
	addr = host[:lastInd]
	port, _ = strconv.Atoi(host[lastInd+1:])
	return
}

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
	dstHost, dstPort := splitAddrPort(r.Host)
	rmIp, rmPort := splitAddrPort(r.RemoteAddr)

	type Request struct {
		Ip      string
		Port    int
		Method  string
		Body    string
		Headers http.Header
	}
	type Response struct {
		Status  int
		Headers http.Header
	}
	type RequestStatus struct {
		Request  Request
		Response Response
		Host     string
		Port     int
		Proto    string
	}
	respJson := RequestStatus{
		Request: Request{
			Method:  r.Method,
			Ip:      rmIp,
			Port:    rmPort,
			Body:    string(body),
			Headers: r.Header,
		},
		Response: Response{
			Status:  httpStatusInt,
			Headers: w.Header(),
		},
		Host:  dstHost,
		Port:  dstPort,
		Proto: r.Proto,
	}
	resp, _ := json.MarshalIndent(respJson, "", "  ")
	_, _ = w.Write(resp)
}

// Error 404 handler
func (s *handler) handlerNotFound(w http.ResponseWriter, r *http.Request) {
	log := s.Service.Log
	errMsg := fmt.Sprintf("failed to get resource: %v", r.RequestURI)
	log.Warn(errMsg)
	jsonError(w, errMsg, 404)
}

// Error 405 handler
func (s *handler) handlerMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	log := s.Service.Log.WithField("service", "rest")
	errMsg := fmt.Sprintf("method %v not allowed", r.Method)
	log.Warn(errMsg)
	jsonError(w, errMsg, 405)
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
