package proxy

import (
	"net/http"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"github.com/maxasm/https-proxy/certs"
	"strings"
	"errors"
)

var dl *log.Logger = log.New(os.Stdout, "[DEBUG] :", log.Lshortfile)

func handle_get_certs(client_info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	// print information about the destination server. SNI
	domain := client_info.ServerName

	// example of SNI is www.google.com
	fmt.Printf("got a HTTPS connection to the server-name: %s\n", domain)

	// check if domain is empty
	if len(domain) == 0 {
		destination := client_info.Conn.RemoteAddr().String()
		dl.Printf("No SNI set. connection to destination: %s\n", destination)
		// get the index of the colon used in the port IP:PORT
		index_of_colon := strings.Index(destination,":")
		if index_of_colon == -1 {
			return nil, errors.New(fmt.Sprintf("the destination address %s has no port.", destination))
		}

		domain = destination[:index_of_colon]
	}	

	// get the certificate for this specific domain
	cert, err__get_cert := certs.Get_certificate_for_domain(domain)
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

	http.Handle("/",http.FileServer(http.Dir("./web")))
	
	dl.Printf("HTTPS server/proxy started on port: %d.", port)
	err__start_server := server.ListenAndServeTLS("","")
	if err__start_server != nil {
		dl.Printf("failed to start HTTPS server on port: %d\n", port)
		os.Exit(1)
	}

	return nil
}
