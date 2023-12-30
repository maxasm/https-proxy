package fserver

import (
	"net/http"
	"fmt"
	"errors"
	"os"
	"golang.org/x/net/websocket"
	"github.com/maxasm/https-proxy/logger"
)

const PORT = 8090

var dl = logger.DL
var wl = logger.WL

var active_ws_connection *websocket.Conn

func Send_WS_Message(payload interface{}) error {
	if active_ws_connection == nil {
		return errors.New("there is not active websocket connection.\n")
	}

	var wc = active_ws_connection 
	// send payload as JSON
	dl.Printf("sending websocket message ...\n")
	return websocket.JSON.Send(wc, &payload)
}

// set the active connection
func handle_websocket(wc *websocket.Conn) {
	dl.Printf("received a websocket connection\n")
	active_ws_connection = wc
	var msg string
	// we need the following loop so that the connection stays up and is not closed
	// we actually dont anticipate any messages from the client
	for {
		err__receive_msg := websocket.Message.Receive(wc, &msg)
		if err__receive_msg != nil {
			wl.Printf("error listening for messages. %s\n", err__receive_msg)
			return
		}
	}
}


// start the static web server to serve the one-page react web interface application
// also handles WebSocket connections to stream data
func Start_server() {
	var dist_dir = "./web-interface/dist"
	var mux = http.NewServeMux()
	// handle serving and routung the SPA react app
	// the current web app is only one page and has no routing, you can therefore use
	// http.FileServer since there is not `react` routing
	mux.Handle("/", http.FileServer(http.Dir(dist_dir)))
	mux.Handle("/ws", websocket.Handler(handle_websocket))
	
	dl.Printf("started file and websocket server on port: %d\n", PORT)
	err__start_server := http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)
	if err__start_server != nil {
		dl.Printf("failed to start server. %s\n", err__start_server)
		os.Exit(1)
	}
}

