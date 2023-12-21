package certs

import (
	"crypto/tls"
	"os"
	"time"
	"io"
	"encoding/pem"
	"errors"
	"crypto/x509"
)

// check if the cert is expired
// returns true if the cert is expired
func cert_expired(domain string) (bool, error) {
	out_dir := certs_dir+domain+"/"
	cert_file, err__open_cert_file := os.OpenFile(out_dir+"fullchain-cert.pem", os.O_RDONLY, 0777)
	if err__open_cert_file != nil {
		return false, err__open_cert_file
	}

	cert_data, err__read_cert_data := io.ReadAll(cert_file)
	if err__read_cert_data != nil {
		return false, err__read_cert_data 
	}

	block, _ := pem.Decode(cert_data)
	if block == nil {
		return false, errors.New("failed to decode ca.pem file\n") 
	} 
	
	cert, err__parse_cert := x509.ParseCertificate(block.Bytes)
	if err__parse_cert != nil {
		return false, err__parse_cert 
	}

	if time.Now().Compare(cert.NotAfter) >= 0 {
		return true, nil
	}
	
	return false, nil
}

func Get_certificate_for_domain(domain string, is_domain_name bool) (*tls.Certificate, error) {
	err__generate_server_cert := generate_server_cert(domain, is_domain_name) 	
	if err__generate_server_cert != nil {
		dl.Printf("failed to create ceritificate for the domain: %s. %s\n", domain, err__generate_server_cert)
		return nil, err__generate_server_cert 
	}
	
	// load the generated certificate
	out_dir := "./.certs/"+domain+"/"
	// TODO: check if the certificate is expired and if so generate new certs for the domain
	// delete certs dir for that given domain and call Get_certitificate_for_domain()
	cert, err__load_cert := tls.LoadX509KeyPair(out_dir+"fullchain-cert.pem", out_dir+"pr-key.pem")
	if err__load_cert != nil {
		dl.Printf("failed to load SSL certificates. %s\n", err__load_cert)
		return nil, err__load_cert
	}

	// get the expiry date of the certificate
	cert_expired, err__check_cert_expired := cert_expired(domain) 
	if err__check_cert_expired != nil {
		return nil, err__check_cert_expired
	}

	if cert_expired {
		// remove the dir containing the certs so that new certs can be generated
		err__remove_dir := os.Remove(out_dir)
		if err__remove_dir != nil {
			dl.Printf("failed to remove dir: %s\n", err__remove_dir)
			return nil, err__remove_dir 
		}
		dl.Printf("certs for domain: %s are expired. Generating new ones ...\n", domain)
		return Get_certificate_for_domain(domain, is_domain_name)
	}

	return &cert, nil
} 



