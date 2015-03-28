package redismq

import (
	"github.com/garyburd/redigo/redis"
)

type MessageQueue struct {
	conn *redis.Conn
}

func NewMessageQueue() *MessageQueue {
	tmpconn, err := redis.Dial("tcp", ":6379")
	return &MessageQueue {
		conn : &tmpconn
	}
}

func (this *MessageQueue) Close() {
	this.conn.Close()
}

func (this *MessageQueue) Enqueue(buffer []byte) error {
	_, err := this.conn.Do("LPUSH", "queue", buffer)
	return err
}

func (this *MessageQueue) GetMessage() ([]byte, error) {
	ret, err := this.conn.Do("BRPOP", "queue")
	return ret.([]byte), err
}