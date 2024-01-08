package main

import (
	"net/http"
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"errors"
	"context"
	"github.com/maxasm/https-proxy/certs"
	"github.com/gorilla/websocket"
)

func handle_get_certs(client_info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	domain := client_info.ServerName
	// dl.Printf("got a HTTPS connection to the server-name: %s\n", domain)

	// is the destination an IP or domain name
	// used to configure subjectAlternativeName to use either IP or DNS
	is_domain_name := true
	// check if domain is empty
	if len(domain) == 0 {
		is_domain_name = false
		destination := client_info.Conn.RemoteAddr().String()
		local_addr := client_info.Conn.LocalAddr().String()
		dl.Printf("No SNI set. RemoteAddr: %s, LocalAddr: %s\n", destination, local_addr)
		// get the index of the colon used in the port IP:PORT
		index_of_colon := strings.Index(destination,":")
		if index_of_colon == -1 {
			return nil, errors.New(fmt.Sprintf("the destination address %s has no port.", destination))
		}

		wl.Printf("You are connecting to a server using an IP. This will not work unless the IP you are connecting to is the IP of host of the proxy.")
		domain = destination[:index_of_colon]
	}	

	// get the certificate for this specific domain
	cert, err__get_cert := certs.Get_certificate_for_domain(domain, is_domain_name)
	if err__get_cert != nil {
		return nil, err__get_cert
	}

	if cert == nil {
		wl.Printf("returned cert is nil. %s\n", err__get_cert)
	}

	return cert, nil
}

func handle_proxy_conn(rw http.ResponseWriter, r *http.Request) {
	// handle websocket connections
	if websocket.IsWebSocketUpgrade(r) {
		var upgrader = websocket.Upgrader{
			// TODO: This is not a permanent fix
			CheckOrigin: func(r *http.Request) bool {return true},
		}
		
		ws_client_conn, err__connect_client := upgrader.Upgrade(rw, r, nil)
		if err__connect_client != nil {
			dl.Printf("failed to upgrade connection to websockets. %s\n", err__connect_client)
		} else {
			dl.Printf("upgraded connection to Websockets.\n")
		}

		// connect to the server
		var path = "wss://"+r.TLS.ServerName+fmt.Sprintf("%s", r.URL)

		var dialer = websocket.Dialer{
			EnableCompression: false,
		}

		// TODO: set up the headers correctly
		// remove repeated headers
		var header = r.Header.Clone()
		header.Del("Sec-Websocket-Version")
		header.Del("Sec-Websocket-Key")
		header.Del("Sec-Websocket-Extensions")
		header.Del("Connection")
		header.Del("Upgrade")
		header.Del("Origin")

		ws_server_conn, _, err__connect_server := dialer.DialContext(context.TODO(), path, header)
		if err__connect_server != nil {
			dl.Printf("failed to connect to server websocket. %s\n", err__connect_server)
		} else {
			dl.Printf("connected to server websocket.\n")
		}

		// read messages from the client and forward them to the server
		go func() {
			for {
				m_type, msg, err__read_msg := ws_client_conn.ReadMessage()
				if err__read_msg != nil {
					wl.Printf("failed to read websocket message from client. %s\n", err__read_msg)
					return
				}

				err__send_msg := ws_server_conn.WriteMessage(m_type, msg)
				if err__send_msg != nil {
					wl.Printf("failed to write message to server. %s\n", err__send_msg)
					return
				}
			}
		}()

		// read the messages from the server and forward them to the client
		for {
			m_type, msg, err__read_msg := ws_server_conn.ReadMessage()
			if err__read_msg != nil {
				wl.Printf("failed to read message from server websocket connection. %s\n", err__read_msg)
				return
			}

			err__send_msg := ws_client_conn.WriteMessage(m_type, msg)
			if err__send_msg != nil {
				wl.Printf("failed to send message to client websocket connection. %s\n", err__send_msg)
				return
			}
		}

		return
	}

	// run the proxy
	server_name := r.TLS.ServerName
	err__intercept := intercept(r, rw, server_name, true)
	if err__intercept != nil {
		wl.Printf("failed to intercept the conncetion to: %s. %s\n", server_name, err__intercept)
	}
}

func start_HTTPS_proxy(port int) error {
	// create the TLS configuration for the HTTPS server
	tls_config := &tls.Config{
		GetCertificate: handle_get_certs,
		GetConfigForClient: func(ch *tls.ClientHelloInfo) (*tls.Config, error) {
			// dl.Printf("** ClientHelloInfo **\n")
			// fmt.Printf("%s\n", *ch)
			return nil, nil
		},
	}
	
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		TLSConfig: tls_config,
	}

	http.HandleFunc("/", handle_proxy_conn)

	dl.Printf("HTTPS server/proxy started on port: %d", port)
	err__start_server := server.ListenAndServeTLS("","")
	if err__start_server != nil {
		dl.Printf("failed to start HTTPS server on port: %d\n", port)
		os.Exit(1)
	}

	return nil
}
