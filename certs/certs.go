package certs

import (
	"os/exec"
	"errors"
	"strings"
	"fmt"
	"os"
	"io"
) 

var paraphrase = "cats"
var certs_dir = "./.certs/"

func run_cmd(prog string, args ...string) error {
	cmd := exec.Command(prog, args...)
	var stderr = strings.Builder{}
	cmd.Stderr = &stderr
	err__run_cmd := cmd.Run()
	stderr_str := stderr.String()
	if err__run_cmd != nil {
		return errors.New(strings.Trim(stderr_str, "\n \t\r")) 
	}
	return nil
}

func generate_ca_private_key() error {
	out_dir := certs_dir+"ca/"
	err__mkdir := os.MkdirAll(out_dir, 0777)
	if err__mkdir != nil {
		return err__mkdir
	}
	
	err__run_cmd := run_cmd(
		"openssl",
		"genrsa",
		"-aes256",
		"-out",
		out_dir+"ca-pkey.pem",
		"-passout",
		"pass:"+paraphrase,
		"2048",
	)
	
	if err__run_cmd != nil {
		return err__run_cmd
	}
	return nil
}

func generate_ca() error {
	out_dir := certs_dir+"ca/"
	err__mkdir := os.MkdirAll(out_dir, 0777)
	if err__mkdir != nil {
		return err__mkdir
	}

	err__generate_ca_pkey := generate_ca_private_key()
	if err__generate_ca_pkey != nil {
		return err__generate_ca_pkey 
	}

	dl.Printf("generated a new CA private key at %sca-pkey.pem\n", out_dir)

	// generate a new CA cert
	err__run_cmd := run_cmd(
		"openssl",
		"req",
		"-new",
		"-x509",
		"-sha256",
		"-days",
		"365",
		"-key",
		out_dir+"ca-pkey.pem",
		"-out",
		out_dir+"ca.pem",
		"-passin",
		"pass:"+paraphrase,
		"-subj",
		"/C=KE/ST=Nairobi/L=Nairobi/O=m9x/CN=m9x.io",
	)

	if err__run_cmd != nil {
		return err__run_cmd
	}

	dl.Printf("generated a new CA cert at %sca.pem", out_dir)
	return nil
} 

func generate_cert_private_key(domain string, id string) error {
	out_dir := certs_dir+domain+"/"+id+"/"

	err__mkdir := os.MkdirAll(out_dir, 0777)
	if err__mkdir != nil {
		return err__mkdir
	}

	err__run_cmd := run_cmd(
		"openssl",
		"genrsa",
		"-out",
		out_dir+"pr-key.pem",
		"2048",
	)

	if err__run_cmd != nil {
		return err__run_cmd
	}
	return nil
}

func generate_csr(domain string, id string) error {
	out_dir := certs_dir+domain+"/"+id+"/"
	
	err__run_cmd := run_cmd(
		"openssl",
		"req",
		"-new",
		"-sha256",
		"-subj",
		"/CN=m9x.io",
		"-key",
		out_dir+"pr-key.pem",
		"-out",
		out_dir+"cert.csr",
	)

	if err__run_cmd != nil {
		return err__run_cmd
	}
	return nil
}

func generate_cert(domain string, id string) error {
	out_dir := certs_dir+domain+"/"+id+"/"
	
	err__run_cmd := run_cmd(
		"openssl",
		"x509",
		"-req",
		"-sha256",
		"-days",
		"365",
		"-in",
		out_dir+"cert.csr",
		"-CA",
		certs_dir+"ca/ca.pem",
		"-CAkey",
		certs_dir+"ca/ca-pkey.pem",
		"-out",
		out_dir+"cert.pem",
		"-extfile",
		out_dir+"config.cnf",
		"-passin",
		"pass:"+paraphrase,
		"-CAcreateserial",
	)

	if err__run_cmd != nil {
		return err__run_cmd
	}

	return nil
}


func generate_config_file(domain string, id string) error {
	config_file := certs_dir+domain+"/"+id+"/config.cnf" 
	f, err__open_file := os.OpenFile(config_file, os.O_RDWR|os.O_CREATE, 0777)
	if err__open_file != nil {
		return err__open_file
	}

	config_text := fmt.Sprintf("subjectAltName=DNS:%s", domain)
	_, err__write_data := f.Write([]byte(config_text))
	if err__write_data != nil {
		return err__write_data
	}

	return nil
}

func make_cert_chain(domain string, id string) error {
	f, err__open_file := os.OpenFile(certs_dir+domain+"/"+id+"/cert.pem", os.O_RDWR, 0777)
	if err__open_file != nil {
		return err__open_file
	}

	cert_data, err__read_data := io.ReadAll(f)
	if err__read_data != nil {
		return err__read_data 
	}

	// read from ca.pem
	caf, err__open_caf := os.OpenFile(certs_dir+"ca/ca.pem", os.O_RDWR, 0777)            
	if err__open_caf != nil {
		return err__open_caf
	}

	ca_data, err__read_ca_data := io.ReadAll(caf)
	if err__read_ca_data != nil {
		return err__read_ca_data
	}

	// create the full-chain file that contains both the ca and the server cert
	fullchain, err__open_fullchain := os.OpenFile(certs_dir+domain+"/"+id+"/fullchain-cert.pem", os.O_RDWR|os.O_CREATE, 0777)
	if err__open_fullchain != nil {
		return err__open_fullchain
	}

	_, err__write1 := fullchain.Write(cert_data)
	if err__write1 != nil {
		return err__write1
	}

	_, err__write2 := fullchain.Write(ca_data)
	if err__write2 != nil {
		return err__write2
	}

	err__close_fullchain := fullchain.Close()
	if err__close_fullchain != nil {
		return err__close_fullchain 
	} 
	 	
	err__close_f := f.Close()
	if err__close_f != nil {
		return err__close_f
	}

	err__close_caf := caf.Close()
	if err__close_caf != nil {
		return err__close_caf
	}

	return nil
}

// created the .cert dir
func init_dir() error {
	err__mkdir := os.Mkdir(".certs", 0777)
	if err__mkdir != nil {
		return err__mkdir
	}
	return nil
}


// the main funtion to generate a server side SSL/TLS certificate.
func generate_server_cert(domain string, id string) error {
	// err__init := init_dir()
	// if err__init != nil {
	// 	return err__init
	// }
	
	err__get_cert_pkey := generate_cert_private_key(domain, id)
	if err__get_cert_pkey != nil {
		return err__get_cert_pkey 
	}
	
	dl.Printf("generated the cert private key for the domain: %s\n", domain)
	
	err__get_csr := generate_csr(domain, id) 
	if err__get_csr != nil {
		return err__get_csr
	}
	
	dl.Printf("generated the csr for the domain: %s\n", domain)

	err__generate_config := generate_config_file(domain, id)
	if err__generate_config != nil {
		return err__generate_config 
	}
	
	dl.Printf("generated config file for domain: %s\n", domain)
	
	err__get_cert := generate_cert(domain, id)
	if err__get_cert != nil {
		return err__get_cert
	}

	dl.Printf("generated the cert for the domain: %s\n", domain)

	err__make_cert_chain := make_cert_chain(domain, id)
	if err__make_cert_chain != nil {
		return err__make_cert_chain
	}

	dl.Printf("generated the certificate chain for the domain: %s\n", domain)
	
	return nil
}