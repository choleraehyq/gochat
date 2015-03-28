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
)

const (
	port = ":1234"
)
var svr server.Server

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	svr = server.NewServer()
}

func clean() {
	//empty now.
}\

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", port)
	checkError(err)
	Listener, err := net.ListenUDP("udp", udpAddr)
	checkError(err)
	
	go svr.Start(Listener)
	
	sigmask()
	
	svr.Stop()
	clean()
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