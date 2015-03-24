package statistic

import (
	"github.com/choleraehyq/gochat/safemap"
	"net"
	"time"
	"sync/atomic"
)

var (
	ConnNum uint32 = 0
	//uid --> conn
	UidConnMap *safemap.Safemap = safemap.Newsafemap()
	//conn --> timestamp
	TimeStampMap *safemap.Safemap = safemap.Newsafemap()
	TryConnect uint32 = 0
)

func RegisterTimeStampConn(conn *net.TCPConn, timeStamp time.Time) {
	TimeStampMap.Set(conn, timeStamp)
}

func UnRegisterTimeStampConn(conn *net.TCPConn) {
	statistic.TimeStampMap.Remove(conn)
}

func AddCount(status, count uint32) {
	atomic.AddUint32(status, count)
}
