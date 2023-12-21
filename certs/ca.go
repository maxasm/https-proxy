package certs

import (
	"os"
	"io"
	"errors"
	"encoding/pem"
	"crypto/x509"
	"time"
)

func generate_ca_cert() error {
	// delete the entire .certs dir
	err__remove_dir := os.Remove(certs_dir)
	if err__remove_dir != nil {
		if !os.IsNotExist(err__remove_dir) { 
			return err__remove_dir 
		}
	}

	// make the .certs dir again
	err__mkdir := os.Mkdir(certs_dir, 0777)
	if err__mkdir != nil {
		return err__mkdir 
	}
	
	// generate a new private key for the CA
	generate_ca_private_key()
	// use the private key to generate a  new CA
	generate_ca()
	dl.Printf("generated a new CA certificate and stored in in .certs/ca.pem\nMake sure to install the CA in you system.\n")
	return nil
	
}

// checks if the CA is expired or if the certs exists
func check_ca() error {
	out_dir := certs_dir+"ca"+"/"
	ca_file, err__open_ca_file := os.OpenFile(out_dir+"ca.pem", os.O_RDWR, 0777)
	if err__open_ca_file != nil {
		if os.IsNotExist(err__open_ca_file) {
			// if the CA ca.pem does not exist, generate a new CA
			return generate_ca_cert()
		}
		return err__open_ca_file
	}

	cert_data, err__read_cert_data := io.ReadAll(ca_file)
	if err__read_cert_data != nil {
		return err__read_cert_data 
	}

	block, _ := pem.Decode(cert_data)
	if block == nil {
		return errors.New("failed to decode ca.pem file\n") 
	} 
	
	ca_cert, err__parse_cert := x509.ParseCertificate(block.Bytes)
	if err__parse_cert != nil {
		return err__parse_cert 
	}

	// check is the CA is expired or the certs dont exist
	if time.Now().Compare(ca_cert.NotAfter) >= 0 {
		dl.Printf("the CA certificate is expired ... generating a new one. This will delete all singed certificates\n")
		return generate_ca_cert()
	}

	dl.Printf("ca certificate is not expired.\n")
	return nil
}


