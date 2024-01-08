package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
	"crypto/rand"
	"encoding/hex"
	"os"
	"sync"
)

// is the web interaface attached?
// if not, do not send the websocket logs
var client_attached bool = true
var rwMutex2 = sync.RWMutex{}

// the active websocket connection that will be used to send the Websocket messages
var active_ws_connection *websocket.Conn = nil

func send_websocket_msg(payload interface{}) error {
	if !client_attached {
		return nil
	}
	if active_ws_connection == nil {
		return nil
	}
	rwMutex2.RLock()
	err__write := active_ws_connection.WriteJSON(payload)
	rwMutex2.RUnlock()
	return err__write
}

// initialize the web interface server
func start_interface_server(port int) {
	dl.Printf("starting HTTP client interface webserver on port: %d\n", port)
	// websocket handler
	var handle_ws_connection = func(rw http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{}
		ws_conn, err__connect := upgrader.Upgrade(rw, r, nil)	
		if err__connect != nil {
			wl.Printf("failed to upgrade websocket connection. %s\n", err__connect) 
			return
		}
		// set the value of the active websocket connection
		rwMutex2.Lock()
		active_ws_connection = ws_conn
		rwMutex2.Unlock()
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./web-interface/dist/")))
	mux.HandleFunc("/ws", handle_ws_connection)
	err__start_server := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err__start_server != nil {
		wl.Printf("failed to start web server on port: %d. %s\n", port, err__start_server)
		os.Exit(1)
	}
}

// TODO: update URL -> Path in client
type RequestInfo struct {
	Method         string              `json:"method,omitempty"`       // the request method
	Headers        map[string][]string `json:"headers,omitempty"`      // the request headers
	URL           string              `json:"url,omitempty"`         // the complete path inluding all URL params
	Payload        string              `json:"payload,omitempty"`            // the payload in the request body (base64 encoded)
	Id             string              `json:"id,omitempty"`           // unique id for this client request
	Response       ResponseInfo        `json:"responseinfo,omitempty"` // the corresponding response to this request
	Protocol       string              `json:"protocol,omitempty"`     // either http/1.1, http/2.0 or http/3.0
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
	rand.Read(buffer)
	return hex.EncodeToString(buffer)
}
