package main

import (
	"github.com/maxasm/https-proxy/proxy"
	"github.com/maxasm/https-proxy/fserver"
)

// TODO:
// 1. handle websocket connections (new feature)
// 2. implement 'open with external application' (new feature) 
// 3. investigate why the speed is slower when using the proxy (debug)
// 4. reduce the number of logs shown on the server to only show important ones (debug)
// 5. fix all certificate issues (debug)

func main() {
	go func(){
		fserver.Start_server()
	}()
	proxy.Start_HTTPS_Proxy(8443)
}
