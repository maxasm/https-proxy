package main

import (
	"github.com/maxasm/https-proxy/proxy"
	"github.com/maxasm/https-proxy/fserver"
)

func main() {
	go func(){
		fserver.Start_server()
	}()
	proxy.Start_HTTPS_Proxy(8443)
}
