package main

import (
	"io"
	"fmt"
	"net"
	"net/http"
	"sync"
	"net/url"
	"bytes"
	"golang.org/x/net/http2"
	"bufio"
	"time"
	tls "github.com/refraction-networking/utls"
)

// For HTTP/2.0 only have one TCP connection per client
var client_connections = make(map[string]*http2.ClientConn)
var rwMutex = sync.RWMutex{}

type BufferCloser struct {
	bf *bytes.Buffer
}

func (bc *BufferCloser) Read(buffer []byte) (int, error) {
	return bc.bf.Read(buffer)
}
// close does nothing
func (bc *BufferCloser) Close() error { return nil }

func intercept(r *http.Request, w http.ResponseWriter, server_name string, is_https bool) error {
	// check what HTTP version is being used
	req_url, err__parse_url := url.Parse("https://"+server_name+fmt.Sprintf("%s", r.URL))
	if err__parse_url != nil {
		return err__parse_url
	}

	// copy all data from the request into a buffer
	var req_buffer = bytes.Buffer{}
	_, err__cp_req := io.Copy(&req_buffer, r.Body)
	if err__cp_req != nil {
		return err__cp_req
	}

	var request_info *RequestInfo = nil
	// send log to client
	if client_attached {
		request_info = &RequestInfo {
			Method: r.Method,
			Headers: r.Header.Clone(),
			URL: fmt.Sprintf("%s", req_url),
			Id: generate_id(),
			Protocol: r.Proto,
			Payload: string(req_buffer.Bytes()),
		}
		err__send_msg := send_websocket_msg(request_info)
		if err__send_msg != nil {
			wl.Printf("failed to send websocket message. %s\n", err__send_msg)
			return err__send_msg
		}
	}
	
	var bf_closer = BufferCloser{bytes.NewBuffer(req_buffer.Bytes())}
	r.Body = &bf_closer
	// update the request url
	(*r).URL = req_url
	resp, err__connect := connect(server_name, r)
	dl.Printf("intercepting connection to SNI %s\n", server_name)
	
	if err__connect != nil {
		wl.Printf("failed to connect to server: %s\n", err__connect)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
		return err__connect
	}

	var resp_buffer = bytes.Buffer{}
	_, err__cp_2 := io.Copy(&resp_buffer, resp.Body)
	if err__cp_2 != nil {
		return err__cp_2
	}
	// close the response body
	resp.Body.Close()
	// send log to client
	if client_attached {
		var response_info = ResponseInfo {
			StatusCode: resp.StatusCode,
			Status: resp.Status, 
			Headers: resp.Header.Clone(),
			Payload: string(resp_buffer.Bytes()),
		}
		request_info.Response = response_info
		err__send_msg := send_websocket_msg(request_info)
		if err__send_msg != nil {
			wl.Printf("failed to send websocket message. %s\n", err__send_msg)
			return err__send_msg
		}
	}

	// copy the response headers
	for header, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(header, value)
		}
	}

	// set the response Status
	w.WriteHeader(resp.StatusCode)
	// read all data from the from the resp to the original response
	_, err__copy := io.Copy(w, &resp_buffer)
	if err__copy != nil {
		return err__copy
	}
	return nil
}

func connect(server_name string, r *http.Request) (*http.Response, error) {
		// check if the cached client connection is still open
		rwMutex.RLock()
		client_conn, exists := client_connections[server_name]
		rwMutex.RUnlock()
		if exists {
			if !(client_conn.State().Closing || client_conn.State().Closed) {
				return handle_http2_conn(client_conn, r)
			}
			delete(client_connections, server_name)
		}

		// create a new client connection
		tls_conn, err__tls_connect := create_tls_conn(server_name)
		if err__tls_connect != nil {
			return nil, err__tls_connect
		}
	
		alpn := tls_conn.ConnectionState().NegotiatedProtocol
		switch alpn {
			case "h2":
				tr := http2.Transport{}
				http2_conn, err__connect := tr.NewClientConn(tls_conn) 
				if err__connect != nil {
					dl.Printf("failed to create new HTTP/2.0 client. %s\n", err__connect)
					return nil, err__connect				
				}				
				rwMutex.Lock()
				client_connections[server_name] = http2_conn
				rwMutex.Unlock()
				return handle_http2_conn(http2_conn, r)
			default: 
				return handle_http1_conn(tls_conn, r)
		}
}


func handle_http1_conn(tls_conn net.Conn, r *http.Request) (*http.Response, error) {
	err__write := r.Write(tls_conn)
	if err__write != nil {
		return nil, err__write
	}
	return http.ReadResponse(bufio.NewReader(tls_conn), r)
}

func handle_http2_conn(http2_conn *http2.ClientConn, r *http.Request) (*http.Response, error) {
	http_resp, err__http2 := http2_conn.RoundTrip(r)
	if err__http2 != nil {
		dl.Printf("failed to make HTTP/2.0 roundTrip. %s\n", err__http2)
	}
	return http_resp, err__http2
}

// create a TLS connection to a server and performs the handshake so as to
// mimic the Firefox broweser. Not sure if this will work for mobile Apps
func create_tls_conn(server_name string) (*tls.UConn, error) {
	// create a TCP connection
	tcp_conn, err__tcp_connect := net.DialTimeout("tcp", server_name+":443", time.Duration(time.Second * 5))
	if err__tcp_connect != nil {
		return nil, err__tcp_connect
	}
	// tls handshake
	tls_config := tls.Config{
		ServerName: server_name,
	}

	tls_conn := tls.UClient(tcp_conn, &tls_config, tls.HelloFirefox_120)

	err__tls_handshake := tls_conn.Handshake()
	if err__tls_handshake != nil {
		return nil, err__tls_handshake 
	}

	return tls_conn, nil
}
