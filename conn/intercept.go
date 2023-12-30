package conn

import (
	"bytes"
	"github.com/maxasm/https-proxy/fserver"
	"github.com/maxasm/https-proxy/logger"
	"io"
	"fmt"
	"net/http"
	"strings"
	"crypto/rand"
	"encoding/hex"
	"encoding/base64"
)

var dl = logger.DL
var wl = logger.WL

type RequestInfo struct {
	Method         string              `json:"method,omitempty"`       // the request method
	Headers        map[string][]string `json:"headers,omitempty"`      // the request headers
	Path           string              `json:"path,omitempty"`         // the complete path inluding all URL params
	Payload        string              `json:"payload,omitempty"`            // the payload in the request body (base64 encoded)
	Id             string              `json:"id,omitempty"`           // unique id for this client request
	Response       ResponseInfo        `json:"responseinfo,omitempty"` // the corresponding response to this request
	Protocol       string              `json:"protocol,omitempty"`     // either http/1.1, http/2.0 or http/3.0
	ConnectionType string              `json:"type,omitempty"`         // either `http` or `web`
}

type ResponseInfo struct {
	StatusCode    int                 `json:"statuscode,omitempty"`    // the response status code
	Status string              `json:"status,omitempty"` // the corresponding response message
	Headers       map[string][]string `json:"headers,omitempty"`       // the response headers
	Payload       string              `json:"payload,omitempty"`       // the response payload
}


// helper function to generate ids
func generate_id() string {
	var buffer = make([]byte, 8)
	_, err__read := rand.Read(buffer)
	if err__read != nil {
		wl.Printf("failed to generate a new ID. %s\n", err__read)
		return "1"
	}
	
	id_str := hex.EncodeToString(buffer)
	return id_str
}

func Intercept(r *http.Request, w http.ResponseWriter, server_name string, is_https bool) error {
	// check what HTTP version is being used
	dl.Printf("intercepting the connection to: %s\n", server_name)

	var full_request_path = "https://"+server_name+fmt.Sprintf("%s", r.URL)
	
	resp,req_info, err__connect := connect(full_request_path, r)
	
	if err__connect != nil {
		wl.Printf("failed to connect to server: %s\n", err__connect)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return err__connect
	}

	// copy the response headers
	for a, b := range resp.Header {
		w.Header().Set(a, strings.Join(b, ","))
	}
	
	// set the response Status
	w.WriteHeader(resp.StatusCode)

	// buffer to read from the response Body -> for logging
	buffer := bytes.Buffer{}

	// read all data from the response into the buffer
	_, err__copy_1 := io.Copy(&buffer, resp.Body)
	if err__copy_1 != nil {
		return err__copy_1
	}

	_, err__copy_2 := io.Copy(w, bytes.NewBuffer(buffer.Bytes()))
	if err__copy_2 != nil {
		return err__copy_2
	}

	// encode the payload using base64
	payload_b64 := base64.StdEncoding.EncodeToString(buffer.Bytes())

	// send the logging data
	var resp_info = ResponseInfo{
		StatusCode: resp.StatusCode, 
		Status: resp.Status,
		Headers: resp.Header.Clone(),
		Payload: payload_b64,
	}

	// set the response value for req_info
	req_info.Response = resp_info

	// send the new response info
	err__send_msg := fserver.Send_WS_Message(req_info)
	if err__send_msg != nil {
		wl.Printf("failed to send websocket message. %s\n", err__send_msg)
	}
	return nil
}

func connect(fpath string, r *http.Request) (*http.Response, *RequestInfo, error) {
	// copy the original request
	req_method := r.Method
	// the reqeust body if the method is a POST request
	buffer := bytes.Buffer{}

	// copy the request body
	_, err__copy_body := io.Copy(&buffer, r.Body)
	if err__copy_body != nil {
		return nil,nil, err__copy_body
	}

	// NOTE: You can set up TLS config here
	client := http.Client{}
	req, err__make_req := http.NewRequest(req_method, fpath, bytes.NewBuffer(buffer.Bytes()))
	if err__make_req != nil {
		return nil,nil, err__make_req
	}

	// TODO: Do headers use `,` or `;` as a delimeter
	// set the headers for the new request
	for a, b := range r.Header {
		req.Header.Set(a, strings.Join(b, ","))
	}

	resp, err__get_resp := client.Do(req)
	if err__get_resp != nil {
		return nil,nil, err__get_resp
	}

	// base64 encode the payload
	payload_b64 := base64.StdEncoding.EncodeToString(buffer.Bytes())

	// create a new ID for the connection
	id := generate_id()

	// send request information for logging
	var req_info = RequestInfo{
		Method: r.Method,
		Path: fpath,
		Headers: r.Header.Clone(),
		ConnectionType: "http",
		Payload: payload_b64,
		Id: id,
		Protocol: r.Proto,
	}

	err__send_msg := fserver.Send_WS_Message(req_info)
	if err__send_msg != nil {
		wl.Printf("failed to send WS message to client. %s\n", err__send_msg)
	}

	return resp, &req_info, nil
}

// func handleWebsocket(wc *websocket.Conn) {
	
// }

// type CustomHandler struct {
// 	cf func(wc *websocket.Conn)
// }

// func (ch CustomHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	// pass it down to websocket.Handler(cf)
// 	if headers == websocket {
// 		websocket.Handler(ch.cf).ServeHttp(w,r)
// 	} else {
// 		// normal Web Handler
// 	}
// } 

// func GetCustomHandler(cf func(wc *websocket.Conn)) Handler {
// 	return CustomHandler{
// 		cf: cf,
// 	}
// }

// mux.Handle("/", GetCustomHandler(handleWebsocket))
