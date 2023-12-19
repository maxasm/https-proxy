package certs

import (
	"crypto/tls"
	"log"
	"os"
)

var dl *log.Logger = log.New(os.Stdout, "[DEBUG]: ", log.Lshortfile)

func Get_certificate_for_domain(domain string) (*tls.Certificate, error) {
	id := "1"
	err__generate_server_cert := generate_server_cert(domain, id) 	
	if err__generate_server_cert != nil {
		dl.Printf("failed to create ceritificate for the domain %s. %s\n", domain, err__generate_server_cert)
		return nil, err__generate_server_cert 
	}
	
	// load the generated certificate
	out_dir := "./.certs/"+domain+"/"+id+"/"
	cert, err__load_cert := tls.LoadX509KeyPair(out_dir+"fullchain-cert.pem", out_dir+"pr-key.pem")
	if err__load_cert != nil {
		dl.Printf("failed to load SSL certificates. %s\n", err__load_cert)
		return nil, err__load_cert
	}

	return &cert, nil
} 


