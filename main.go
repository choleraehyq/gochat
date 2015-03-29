// gochat project main.go
package main

import (
	"github.com/choleraehyq/gochat/server"
	"fmt"
	"net"
	"runtime"
	"os"
	"os/signal"
	"syscall"
	"log"
	"flag"
)

var svr server.Server
var kind, port, serverAddr string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	
	kind = flag.String("kind", "server", "client for client and server for server")
	port = flag.String("port", ":1234", "udp listening port")
	serverAddr = flag.String("serverAddr", "localhost:1234", "address and port of the server")
	flag.Parse()
}

func clean() {
	//empty now.
}\

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	checkError(err)
	Listener, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	
	if kind == "server" {
		svr = server.NewServer()
		go svr.Start(Listener)
	
		sigmask()
	
		svr.Stop()
		clean()
	}
	else if kind == "client" {
		cli := client.NewClient(Listener)
		serverUdp , err := net.ResolveUDPAddr("udp", serverAddr)
		checkError(err)
		go cli.Start(serverUdp)
		
		sigmask()
		cli.Stop()
		clean()
	}
	else {
		log.Println("cli parameter wrong")
		os.Exit(1)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error: %v\n", err)
		os.Exit(1)
	}
}

func sigmask() {
	sigCh := make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("Signal: %v\n", <-sigCh)	
}