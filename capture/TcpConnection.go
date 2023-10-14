package capture

import "time"

const (
	ack = iota
	ackOk
)

type TCPConnection struct {
	LastTime time.Time
	payLoad  []byte
	dataChan chan<- []byte
}

func (tcp *TCPConnection) AaddPacket(payload []byte, lastTime time.Time) {
	tcp.payLoad = append(tcp.payLoad, payload...)
	tcp.LastTime = lastTime
}

func (tcp *TCPConnection) ResultData(){
	tcp.dataChan<-tcp.payLoad
}

func NewTCPConnection(lastTime time.Time,dataChan chan<- []byte) *TCPConnection {
	return &TCPConnection{
		LastTime: lastTime,
		payLoad:  make([]byte, 0),
		dataChan: dataChan,
	}
}
