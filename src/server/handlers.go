package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func splitAddrPort(host string) (addr string, port int) {
	lastInd := strings.LastIndex(host, ":")
	addr = host[:lastInd]
	port, _ = strconv.Atoi(host[lastInd+1:])
	return
}

// App index page
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
	body, _ := ioutil.ReadAll(r.Body)
	dstHost, dstPort := splitAddrPort(r.Host)
	rmIp, rmPort := splitAddrPort(r.RemoteAddr)
	respJson := struct {
		RmIp    string
		RmPort  int
		DstHost string
		DstPort int
		Method  string
		Body    string
		Headers http.Header
		Status  int
	}{
		RmIp:    rmIp,
		RmPort:  rmPort,
		DstHost: dstHost,
		DstPort: dstPort,
		Method:  r.Method,
		Body:    string(body),
		Headers: r.Header,
		Status:  httpStatusInt,
	}
	resp, err := json.MarshalIndent(respJson, "", "  ")
	if err != nil {
		fmt.Println("failed to marshal json")
	}
	w.WriteHeader(httpStatusInt)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("error index page")
	}
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

// Return json-formatted HTTP error
func jsonError(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json")
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
