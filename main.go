package main

// TODO:
// 0. Investigate why I am getting the forbidden error. This is what is causing the CORS errors + other errors
// 1. handle websocket connections (new feature) + instagram bad handshake (debug)
// 2. can't log into netflix (debug)
// 3. Some times there is no response code, why is this? (bug) + websocket logging if web interface client is not attached
// 4. investigate why the speed is slower when using the proxy (debug)
// 5. reduce the number of logs shown on the server to only show important ones (debug)
// 6. fix all certificate issues (debug)
// 7. WhatsApp not loading all resources as expected
// 8. Fix curl issues -> Research on how to properly formart the ceritificate. (debug)
// 9. Fix panic when there is no connection to the internet.


func main() {
	if client_attached {
		go func() {
			start_interface_server(8444)
		}()
	}
	start_HTTPS_proxy(8443)
}
