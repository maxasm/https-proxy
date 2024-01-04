package conn

import (
	"bytes"
	"github.com/maxasm/https-proxy/fserver"
	"github.com/maxasm/https-proxy/logger"
	"io"
	"os"
	"fmt"
	"net/http"
	"crypto/rand"
	"encoding/hex"
	"encoding/base64"
	"compress/gzip"
	"compress/flate"
	"github.com/andybalholm/brotli"
)

var dl = logger.DL
var wl = logger.WL

// TODO: rename in JSON URL -> path
type RequestInfo struct {
	Method         string              `json:"method,omitempty"`       // the request method
	Headers        map[string][]string `json:"headers,omitempty"`      // the request headers
	URL           string              `json:"path,omitempty"`         // the complete path inluding all URL params
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

func decode_br(enc_data []byte) []byte {
	var data_reader = bytes.NewReader(enc_data)
	var deflate_reader = brotli.NewReader(data_reader)
	
	var decoded_data = bytes.Buffer{}
	_, err__copy := io.Copy(&decoded_data, deflate_reader)
	if err__copy != nil {
		wl.Printf("failed to copy data from brtoli decoded data. %s\n", err__copy)
		return enc_data
	}
	return decoded_data.Bytes()
}

func decode_gzip(enc_data []byte) []byte {
	var data_reader = bytes.NewReader(enc_data)
	var gzip_reader, err__gzip_reader = gzip.NewReader(data_reader)
	if err__gzip_reader != nil {
		wl.Printf("failed to create a new gzip Reader. %s\n", err__gzip_reader)
		return enc_data
	}
	var decoded_data = bytes.Buffer{}
	_, err__copy := io.Copy(&decoded_data, gzip_reader)
	if err__copy != nil {
		wl.Printf("failed to copy gzip decoded data. %s\n", err__copy)
		return enc_data
	}
	return decoded_data.Bytes()
}

func decode_flate(enc_data []byte) []byte {
	var data_reader = bytes.NewReader(enc_data)
	var deflate_reader = flate.NewReader(data_reader)
	var decoded_data = bytes.Buffer{}
	_, err__copy := io.Copy(&decoded_data, deflate_reader)
	if err__copy != nil {
		wl.Printf("failed to copy data from flate decoded data. %s\n", err__copy)
		return enc_data
	}
	return decoded_data.Bytes()
}

func Intercept(r *http.Request, w http.ResponseWriter, server_name string, is_https bool) error {
	// check what HTTP version is being used
	dl.Printf("intercepting the connection to: %s\n", server_name)
	var url = "https://"+server_name+fmt.Sprintf("%s", r.URL)
	resp,req_info, err__connect := connect(url, r)
	
	if err__connect != nil {
		wl.Printf("failed to connect to server: %s\n", err__connect)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return err__connect
	}

	// copy the response headers
	for header, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	// set the response Status
	w.WriteHeader(resp.StatusCode)

	// buffer to read from the response Body
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

	// decode the data based on the Content-Encoding header
	var payload_data = buffer.Bytes()
	var content_encoding_alg = resp.Header.Values("Content-Encoding")
	if len(content_encoding_alg) != 0 {
		if len(content_encoding_alg) > 1 {
			// TODO: Handle more than one encoding schemes
			wl.Printf("more than one encoding schemes using used.\n")
			os.Exit(1)
		}

		var encoding_alg = content_encoding_alg[0] 
		switch encoding_alg {
			case "br":
				payload_data = decode_br(payload_data)
			case "gzip":
				payload_data = decode_gzip(payload_data)
			case "deflate":
				payload_data = decode_flate(payload_data)
			default:
				wl.Printf("Encoding %s not supported. %s\n", encoding_alg)
		}
	}
	// encode the payload using base64
	var payload_b64 = base64.StdEncoding.EncodeToString(payload_data)

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

func connect(url string, r *http.Request) (*http.Response, *RequestInfo, error) {
	buffer := bytes.Buffer{}
	// copy the request body
	_, err__copy_body := io.Copy(&buffer, r.Body)
	if err__copy_body != nil {
		return nil,nil, err__copy_body
	}

	// NOTE: You can set up TLS config here
	client := http.Client{}
	req, err__make_req := http.NewRequest(r.Method, url, bytes.NewReader(buffer.Bytes()))
	if err__make_req != nil {
		return nil,nil, err__make_req
	}

	// set the headers for the new request
	for header, values := range r.Header {
		for _, value := range values {
			req.Header.Add(header, value)
		}
	}

	resp, err__get_resp := client.Do(req)
	if err__get_resp != nil {
		return nil,nil, err__get_resp
	}

	// base64 encode the payload
	payload_b64 := base64.StdEncoding.EncodeToString(buffer.Bytes())

	// create a new ID for the connection
	id := generate_id()
	
	dl.Printf("Access-Control-Allow-Origin header for %s is: %s\n", url, resp.Header.Get("Access-Control-Allow-Origin"))
	// send request information for logging
	var req_info = RequestInfo{
		Method: r.Method,
		URL: url,
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
