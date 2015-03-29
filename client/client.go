package client 

import (
	"os"
	"code.google.com/p/goprotobuf/proto"
	"github.com/choleraehyq/gochat/msgproto"
	"fmt"
	"log"
	"net"
	"time"
)

type Client struct {
	conn *net.UDPConn
	server *net.UDPAddr
	exitCh chan bool
}

func NewClient(listener *net.UDPConn, addr *net.UDPAddr) *Client {
	return &Client {
		conn: listener,
		server: addr,
	}
} 

func (this *Client) connectToServer() {
	msg := &msgproto.SendMsg{
		MsgType: proto.String("Register")
	}
	buffer, err := proto.Marshal(msg)
	checkError(err)
	_, err := this.conn.WriteTo(buffer, this.server)
	checkError(err)
}

func (this *Client) receiveMessage() bool {
	this.conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	this.conn.SetReadBuffer(1024)
	var buffer []byte
	_, _, err := this.conn.ReadFromUDP(buffer)
	if isTimeout(err) {
		return false
	}
	msg := &msgproto.SendMsg{}
	err := proto.Unmarshal(buffer, msg)
	if err != nil {
		log.Printf("Encoding error: %v\n", err)
		return false
	}
	fmt.Printf("Receive from %v\n", msg.GetMsgFrom())
	fmt.Printf("Content is: %v\n", msg.GetMsgContent())
	return true
}

func (this *Client) sendMessage() bool {
	fmt.Println("Enter the target address")
	var addrString string
	fmt.Scanf("%s", addrString)
	var content string
	fmt.Scanf("%s", content)
	msg := &msgproto.SendMsg{
		MsgFrom: proto.String(this.conn.LocalAddr().String())
		MsgSendTo: proto.String(addrString)
		MsgType: proto.String("Send")
		MsgContent: proto.String(content)
	}
	buffer, err := proto.Marshal(msg)
	checkError(err)
	_, err := this.conn.WriteTo(buffer, this.server)
	checkError(err)
	fmt.Println("Exit?[Y/N]")
	var isExit string
	fmt.Scanf("%s", isExit)
	if isExit == "Y" {
		return false
	}
	else {
		return true
	}
}

func (this *Client) Start() {
	cli.connectToServer(serverUdp)
	for {
		this.exitOrNot()
		this.receiveMessage()
		fmt.Println("Do you want to send message?[Y/n]")
		var choose string
		fmt.Scanf("%s", choose)
		if choose == "Y" {
			while (this.sendMessage()) {
				this.receiveMessage()
			}
		}
	}
}

func (this *Client) exitOrNot() {
	select {
		case <-this.exitCh:
			log.Printf("Client stop\n")
			return
		default:
	}
}

func isTimeout(err error) bool {
	e, ok := err.(Error)
	return ok && e.Timeout()
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("Error: %v\n", err)
		os.Exit(1)
	}
}