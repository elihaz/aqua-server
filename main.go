package main

import (
	"log"
	"fmt"
	"flag"
)

var protocol string
var format string
var port int

func init() {

	flag.StringVar(&protocol, "protocol", "HTTP", "the protocol type of the server, accepts HTTPS, HTTP")
	flag.StringVar(&format, "format", "JSON", "The format of file requests, accepts XML, JSON")
	flag.IntVar(&port, "port", 5000, "HTTP service address")
	flag.Parse()

}

func main() {

	if format != "JSON" && format != "XML" {
		log.Print("Supported formats are JSON and XML only")
		return
	}

	if protocol != "HTTP" && protocol != "HTTPS" {
		log.Print("Supported protocols are HTTP and HTTPS only")
		return
	}

	aquaServer := newServer(port, protocol, format)
	log.Print(fmt.Sprintf("Started %s server Listening on port %d", protocol, port))
	aquaServer.start()
}
