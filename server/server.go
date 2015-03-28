package server

import (
	"github.com/choleraehyq/gochat/statistic"
	"github.com/choleraehyq/gochat/redismq"
	"fmt"
	"sync"
	"net"
	"time"
	"log"
	"os"
)

type Server struct {
	exitCh chan bool
	waitGroup *sync.WaitGroup
	maxPacLen uint32
	mq *redismq.MessageQueue
	acceptTimeout time.Duration    
	readTimeout time.Duration    
	writeTimeout time.Duration  
	spamTimeout time.Duration
}

func NewServer() *Server {
	return &Server {
		exitCh: make(chan bool),
		waitGroup: &sync.WaitGroup{},
		maxPacLen: 2048,
		mq: redismq.NewMessageQueue,
		acceptTimeout: 60,
		readTimeout:   60,
		writeTimeout:  60,
	}
}


func (this *Server) SetAcceptTimeout(acceptTimeout time.Duration) {
	this.acceptTimeout = acceptTimeout
}

func (this *Server) SetReadTimeout(readTimeout time.Duration) {
	this.readTimeout = readTimeout
}

func (this *Server) SetWriteTimeout(writeTimeout time.Duration) {
	this.writeTimeout = writeTimeout
}

func (this *Server) SetSpamTimeout(spamTimeout time.Duration) {
	this.spamTimeout = spamTimeout
}

func (this *Server) SetMaxPacLen(maxPacLen uint32) {
	this.maxPacLen = maxPacLen
}

func (this *Server) Start(listener *net.UDPConn) {
	log.Printf("Server start at %v\n", listener.Addr())
	this.waitGroup.Add(1)
	defer func(){
		listener.Close()
		this.waitGroup.Done()
	}()
	
	listener.SetReadBuffer(this.maxPacLen)
	go this.dealWithSpamConn()
	
	go this.sendMessage()
	
	for {
		this.exitOrNot(listener)
		
		var buffer []byte
		_, addr, err := listener.ReadFromUDP(buffer)
		
		
		statistic.RegisterTimeStampAddr(addr, time.Now())
		log.Printf("Accept successfully from: %v\n", addr)
		go this.handleClientConn(buffer)
	}	
}

func (this *Server) dealWithSpamConn() {
	ticker := time.NewTicker(this.spamTimeout)
	for _ = range ticker.C {
		items := statistic.TimeStampMap.Items()
		for conn, timeStamp := range items {
			if timeStamp != nil {
				deadLine := timeStamp.(time.Time).Add(this.spamTimeout)
				if time.Now().After(deadLine) {
					statistic.UnRegisterTimeStampConn(conn)
				}
			}
		}
	}
}

func (this *Server) exitOrNot(listener *net.TCPListener) {
	select {
		case <-this.exitCh:
			log.Printf("Server stop on %v\n", listener)
			return
		default:
	}
}

func (this *Server) handleClientConn(buffer []byte) {
	//protobuf decode
	//push in mq 
}

func (this *Server) sendMessage() {
	
}

func (this *Server) Stop() {
	close(this.exitCh)
	this.waitGroup.Wait()
}

func checkError(err error) {
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}