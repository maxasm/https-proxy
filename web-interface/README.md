# Web Interface

- A simple Web Interface for searching through Web HTTP Traffic

## Feature List

- A List of all incoming requests from the client.
- The client entry should contain.
  - The protocol [HTTP/1.1 HTTP/2.0 HTTP/3]
  - The request method (GET, POST, PUT, ...) Special symbol for (Websockets Handshake)
  - The response status or code (Pending, 200 OK, 404 NOT FOUND ...)
  - Response content-type and direct access to it ... preview or smthng.
  - server-name

GET | HTTP/1.1   | www.google.com | image/jpeg | 200 | 200KB 
    | WS  | www.google.com | message    |        | 12 B
POST | HTTP/2.0   | www.youtube.com | image/jpeg | 404 | 89 KB 
GET | HTTP/2.0   | www.youtube.com | image/jpeg | 404 | 89 KB 

In depth features:
- Break down a URL based on the parameters in it
- Have a table for cookies
- Easy copy and paste for headers
- An easy way to download the content payload
- An easy way to preview content payload. Donwload the content and open it using an external program sucha as VLC

For Web-Socket, get the current message being sent and all previous messages.

## Query Language
rq.MT == GET & rq.SN == "www.google.com" & rs.Content-Type == "image/jpeg"
rq.PT == WS & rq.MSG == "match" 
rs.SC != 200 & rs.SC != 400 & (rs.HD["Content-Length"] >= 400)

type Request struct {
  Method string // -> req.Method
}

type Response struct {
  
}

Opearators
& and operator
| or operator
() grouping operator
== 'equal to' *Matches*
!= 'not equal to' *Doesn't Match*
>=
<=

rq.MT == "GET"

-> WebSockets -> JSON feed.


function filter(list, args ... ) {
  
  return list
}
