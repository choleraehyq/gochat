package statistic

import (
	"github.com/choleraehyq/gochat/utils/safemap"
	"net"
	"time"
	"sync/atomic"
)

var (
	//conn --> timestamp
	TimeStampMap *safemap.Safemap = safemap.Newsafemap()
	NameConnMap *safemap.Safemap = safemap.Newsafemap()
	uint32 PacketNum
	uint32 BadPacketNum
)

func RegisterTimeStampAddr(conn *net.UDPAddr, timeStamp time.Time) {
	TimeStampMap.Pur(conn, timeStamp)
}

func UnRegisterTimeStampAddr(conn *net.UDPAddr) {
	statistic.TimeStampMap.Remove(conn)
}

func UpdateTimeStampAddr(conn *net.UDPAddr, timeStamp time.Time) {
	TimeStampMap.Put(conn, timeStamp)
}

func AddCount(status, count uint32) {
	atomic.AddUint32(status, count)
}
