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

var svr server.Server

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	svr = server.NewServer()
	
}

func clean() {
	//empty now.
}\

func main() {
	Listener, err := net.ListenTCP("tcp", "127.0.0.1:1234")
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