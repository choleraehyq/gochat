package server

import (
	"github.com/choleraehyq/gochat/statistic"
	"fmt"
	"sync"
	"net"
	"time"
	"log"
)

type Server struct {
	exitCh chan bool
	waitGroup *sync.WaitGroup
	maxPacLen uint32
	acceptTimeout time.Duration    
	readTimeout time.Duration    
	writeTimeout time.Duration  
	spamTimeout time.Duration
}

func NewServer() *Server {
	return &Server {
		exitCh: make(chan bool),
		waitGroup: &sync.WaitGroup{},
		maxPacLen: 256,
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

func (this *Server) Start(listener *net.TCPListener) {
	log.Printf("Server start at %v\n", listener.Addr())
	this.waitGroup.Add(1)
	defer func(){
		listener.Close()
		this.waitGroup.Done()
	}()
	
	go this.dealWithSpamConn()
	
	for {
		this.exitOrNot(listener)
		listener.SetDeadline(time.Now().Add(this.acceptTimeout))
		
		conn, err := listener.Accept()
		if err != nil {
			statistic.AddCount(statistic.TryConnect, 1)
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			log.Printf("Accept error: %v\n", err)
			continue
		}
		
		statistic.AddCount(statistic.ConnNum, 1)
		statistic.RegisterConn(conn, time.Now())
		log.Printf("Accept successfully: %v\n", conn.RemoteAddr())
		go this.handleClientConn(conn)
	}	
}

func (this *Server) dealWithSpamConn() {
	ticker := time.NewTicker(this.spamTimeout)
	for _ = range ticker.C {
		items := statistic.TimeStampMap.Items()
		for conn, timeStamp := range items {
			if timeStamp != nil {
				deadLine := timeStamp.(time.Time)Add(this.spamTimeout)
				if time.Now().After(deadLine) {
					statistic.UnRegisterTimeStampConn(conn)
					conn.(*net.TCPConn).Close()
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

func (this *Server) handleClientConn(conn *net.TCPConn) {
	
}

func (this *Server) Stop() {
	close(this.exitCh)
	this.waitGroup.Wait()
}