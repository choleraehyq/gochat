package server

import (
	"github.com/choleraehyq/gochat/statistic"
	"github.com/choleraehyq/gochat/redismq"
	"code.google.com/p/goprotobuf/proto"
	"github.com/choleraehyq/gochat/msgproto"
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
	//all the Timeout is ms
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
		spamTimeout: 30,
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
	
	go this.sendMessage(listener)
	
	for {
		this.exitOrNot(listener)
		
		var buffer []byte
		listener.SetReadDeadline(time.Now().Add(this.readTimeout * time.Millisecond))
		_, addr, err := listener.ReadFromUDP(buffer)
		if isTimeout(err) {
			continue
		}
		
		this.handleClientConn(buffer)
	}	
}

func (this *Server) dealWithSpamConn() {
	ticker := time.NewTicker(this.spamTimeout)
	for _ = range ticker.C {
		items := statistic.TimeStampMap.Items()
		for conn, timeStamp := range items {
			if timeStamp != nil {
				deadLine := timeStamp.(time.Time).Add(this.spamTimeout * time.Minute)
				if time.Now().After(deadLine) {
					statistic.UnRegisterTimeStampConn(conn)
				}
			}
		}
	}
}

func (this *Server) handleClientConn(buffer []byte) {
	msg := &msgproto.SendMsg{}
	err := proto.Unmarshal(buffer, msg)
	if err != nil {
		log.Printf("Decoding Error: %v\n", err)
	}
	switch msgType := msg.GetMsgType(); msgType {
	case "Register":
		this.onRegister(msg.GetMsgFrom())
		
	//TODO: ping
	//case "Ping":
	//	this.onPing(msg.GetMsgFrom(), msg.GetMsgSendTo())
	
	case "Send":
		this.onSend(buffer)
	default:
		//ignore
	} 
}

func (this *Server) sendMessage(listener *net.UDPConn) {
	this.waitGroup.Add()
	defer this.waitGroup.Done()
	for {
		this.exitOrNot()
		buf, err := this.mq.GetMessage()
		checkError()
		if buf == nil {
			continue
		}
		msg := &msgproto.SendMsg{}
		err := proto.Unmarshal(buf, msg)
		addr := msg.GetMsgSendTo()
		_, err := listener.WriteTo(buf, addr)
		if err != nil {
			log.Printf("Send Error: %v\n", err)
		}
	}
}

func (this *Server) Stop() {
	close(this.exitCh)
	this.waitGroup.Wait()
}

func (this *Server) exitOrNot(listener *net.TCPListener) {
	select {
		case <-this.exitCh:
			log.Printf("Server stop on %v\n", listener)
			return
		default:
	}
}

func (this *Server) onRegister(addr string) {
	statistic.RegisterTimeStampAddr(addr, time.Now())
	log.Printf("Accept: %v\n", addr)
}

func (this *Server) onSend([]byte buffer) {
	this.mq.Enqueue(buffer)
}

func isTimeout(err error) bool {
	e, ok := err.(Error)
	return ok && e.Timeout()
}

func checkError(err error) {
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}