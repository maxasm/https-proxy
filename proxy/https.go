package proxy

import (
	"net/http"
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"errors"
	"github.com/maxasm/https-proxy/certs"
	"github.com/maxasm/https-proxy/logger"
)


var dl = logger.DL
var wl = logger.WL

func handle_get_certs(client_info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	// print information about the destination server. SNI
	domain := client_info.ServerName

	// example of SNI is www.google.com
	dl.Printf("got a HTTPS connection to the server-name: %s\n", domain)

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

	return cert, nil
}

func Start_HTTPS_Proxy(port int) error {
	// create the TLS configuration for the HTTPS server
	tls_config := &tls.Config{
		GetCertificate: handle_get_certs,
	}
	
	server := http.Server{
		Addr: fmt.Sprintf(":%d", port),
		TLSConfig: tls_config,
	}

	// TODO:
	// this is the entry point for all connections
	// use connection infor such as domain, headers, cookies, etc
	// to recreate the original HTTPS connection and send it to the server.
	http.Handle("/",http.FileServer(http.Dir("./web")))
	
	dl.Printf("HTTPS server/proxy started on port: %d.", port)
	err__start_server := server.ListenAndServeTLS("","")
	if err__start_server != nil {
		dl.Printf("failed to start HTTPS server on port: %d\n", port)
		os.Exit(1)
	}

	return nil
}
