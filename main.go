package main

import (
	"os"
	"log"
)

var dl *log.Logger = log.New(os.Stdout, "[DEBUG] :", log.Lshortfile)

func main() {
	// err := generate_ca()
	// if err != nil {
	// 	dl.Printf("error getting ca: %s\n", err)
	// }

	err__generate_cert := generate_server_cert("google.com")
	if err__generate_cert != nil {
		dl.Printf("failed to generated cert. %s\n", err__generate_cert)
		os.Exit(1)
	}
}
