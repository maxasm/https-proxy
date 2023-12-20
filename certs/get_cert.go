package certs

import (
	"crypto/tls"
)

func Get_certificate_for_domain(domain string, is_domain_name bool) (*tls.Certificate, error) {
	err__generate_server_cert := generate_server_cert(domain, is_domain_name) 	
	if err__generate_server_cert != nil {
		dl.Printf("failed to create ceritificate for the domain %s. %s\n", domain, err__generate_server_cert)
		return nil, err__generate_server_cert 
	}
	
	// load the generated certificate
	out_dir := "./.certs/"+domain+"/"
	cert, err__load_cert := tls.LoadX509KeyPair(out_dir+"fullchain-cert.pem", out_dir+"pr-key.pem")
	if err__load_cert != nil {
		dl.Printf("failed to load SSL certificates. %s\n", err__load_cert)
		return nil, err__load_cert
	}

	return &cert, nil
} 


