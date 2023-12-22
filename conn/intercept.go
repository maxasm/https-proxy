package conn

import (
	"net/http"
	"bytes"
	"strings"
	"io"
	"github.com/maxasm/https-proxy/logger"
)

var dl = logger.DL
var wl = logger.WL

func Intercept(r *http.Request, w http.ResponseWriter, server_name string, is_https bool) error {
	// check what HTTP version is being used
	dl.Printf("intercepting the connection to: %s\n", server_name)

	var full_request_path string
	request_path := r.URL.Path
	// the request path may not fail to include the leading slash `/`
	if request_path[0] != '/' {
		full_request_path = "https://"+server_name+"/"+request_path
	} else {
		full_request_path = "https://"+server_name+request_path
	}

	resp, err__connect := connect(full_request_path, r)
	if err__connect != nil {
		wl.Printf("failed to connect to server: %s\n", err__connect)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return err__connect
	}

	// copy the response headers
	for a,b := range resp.Header {
		w.Header().Set(a, strings.Join(b, ","))
	}

	dl.Printf("----- set response headers as ------\n")
	for a, b := range w.Header() {
		dl.Printf("%s: %s\n", a, b)
	}
	dl.Printf("------\n")
	
	// set the response Status
	w.WriteHeader(resp.StatusCode)

	// read all data from the response
	_, err__copy_body := io.Copy(w, resp.Body)
	if err__copy_body != nil {
		wl.Printf("error -> %s\n", err__copy_body)
		return err__copy_body
	}

	return nil
}

func connect(fpath string, r *http.Request) (*http.Response, error) {
	// copy the original request 
	req_method := r.Method
	// the reqeust body if the method is a POST request
	buffer := bytes.Buffer{}

	// copy the request body
	_, err__copy_body := io.Copy(&buffer, r.Body)
	if err__copy_body != nil {
		return nil, err__copy_body
	}

	dl.Printf("connecting ... %s -> %s -> Buffer [%s]\n", req_method, fpath, buffer.String())

	// TODO: You can set up TLS config here
	client := http.Client{}
	req, err__make_req := http.NewRequest(req_method, fpath, &buffer)
	if err__make_req != nil {
		return nil, err__make_req
	}

	// set the headers for the new request
	for a,b := range r.Header {
		req.Header.Set(a, strings.Join(b, ","))
	}

	dl.Printf("----- set the request headers as -----\n")
	for a,b := range req.Header {
		dl.Printf("%s: %s\n", a, b)
	}
	dl.Printf("-----\n")
	resp, err__get_resp := client.Do(req)
	if err__get_resp != nil {
		return nil, err__get_resp
	}

	return resp, nil
}
