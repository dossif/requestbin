package server

import (
	"encoding/json"
	"fmt"
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
func (s *handler) handlerIndex(w http.ResponseWriter, r *http.Request) {
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
	}{
		RmIp:    rmIp,
		RmPort:  rmPort,
		DstHost: dstHost,
		DstPort: dstPort,
		Method:  r.Method,
		Body:    string(body),
		Headers: r.Header,
	}
	resp, err := json.MarshalIndent(respJson, "", "  ")
	if err != nil {
		fmt.Println("failed to marshal json")
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		fmt.Println("error index page")
	}
}

// Error 404 handler
func (s *handler) handlerNotFound(w http.ResponseWriter, r *http.Request) {
	log := s.Service.Log.WithField("service", "rest")
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
