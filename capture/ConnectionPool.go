package capture

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/ip4defrag"
	"github.com/google/gopacket/layers"
	"log"
	"time"
)

type IPRefragKey struct {
	Id    uint16
	SrcIP [4]byte
}

type ConnMsg struct {
	srcIP, dstIP [4]byte
	Start        time.Time
	Last         time.Time
}

type ConnectionPool struct {
	FragmentList *ip4defrag.IPv4Defragmenter
	connMsgs     map[IPRefragKey]*ConnMsg
	TCPList      map[string]*TCPConnection
	mapKeyQueue  *Queue
	DataChan chan<- []byte
}

func (pool *ConnectionPool) DisposePacket(packet gopacket.Packet) {
	//判断流量包网络层的类型，如果为ipv4，则继续执行
	ipv4 := packet.Layer(layers.LayerTypeIPv4)
	if ipv4 == nil {
		//log.Println("(Sniffer.disposePacket) packet parsing fail!")
		return
	}

	//包的类型转换
	ipv4Layer, ok := ipv4.(*layers.IPv4)
	if !ok {
		log.Println("cast failed!")
		return
	}

	var srcIp, dstIp [4]byte
	copy(srcIp[:4], ipv4Layer.SrcIP)
	copy(dstIp[:4], ipv4Layer.DstIP)

	//记录每个IP的时间，主要用于IP分片重组，记录第一个分片和最后一个分配到达时间
	//通过IP数据包中，每个分片的ID是唯一标识的
	ipRefraKey := IPRefragKey{
		Id:    ipv4Layer.Id,
		SrcIP: srcIp,
	}

	t, ok := pool.connMsgs[ipRefraKey]
	if ok {
		t.Last = packet.Metadata().Timestamp
	} else {
		now := &ConnMsg{
			srcIP: srcIp,
			dstIP: dstIp,
			Start: packet.Metadata().Timestamp,
			Last:  packet.Metadata().Timestamp,
		}
		pool.connMsgs[ipRefraKey] = now
	}

	//检测是否需要ip分片重组
	ipPacket, err := pool.FragmentList.DefragIPv4WithTimestamp(
		ipv4Layer, packet.Metadata().Timestamp)
	if err != nil {
		// log.Println("该包为IP分片！(ConversationPool.go 84)")
		return
	}

	pool.checkTimeout(packet.Metadata().Timestamp)

	if ipPacket==nil{
		return
	}
	payload := ipPacket.Payload

	switch ipPacket.Protocol {
	case layers.IPProtocolTCP:
		p := gopacket.NewPacket(payload, layers.LayerTypeTCP, gopacket.Default)
		if layer := p.Layer(layers.LayerTypeTCP); layer != nil {
			if tcp, ok := layer.(*layers.TCP); ok {
				//log.Println(tcp.Payload)
				pool.disposeTCPPackets(tcp, pool.connMsgs[ipRefraKey])
			}
		}

	default:
		//log.Println("未知的包类型 ConversationPool.go ", ipPacket.Protocol)
	}

	delete(pool.connMsgs, ipRefraKey)
}

func (pool *ConnectionPool) disposeTCPPackets(tcp *layers.TCP, connmsg *ConnMsg) {

	fiveTuple := FiveTuple{
		SrcIP:        connmsg.srcIP,
		DstIP:        connmsg.dstIP,
		SrcPort:      uint16(tcp.SrcPort),
		DstPort:      uint16(tcp.DstPort),
		ProtocolType: layers.IPProtocolTCP,
	}

	connectionHash := fiveTuple.Hash()

	if connection, ok := pool.TCPList[connectionHash]; ok {
		pool.mapKeyQueue.ResetValue(connectionHash)
		connection.AaddPacket(tcp.Payload, connmsg.Last)
	} else {

		newConnection := NewTCPConnection(connmsg.Last,pool.DataChan)
		newConnection.AaddPacket(tcp.Payload, connmsg.Last)

		pool.mapKeyQueue.Push(connectionHash)
		pool.TCPList[connectionHash] = newConnection
	}
}

func (pool *ConnectionPool) checkTimeout(now time.Time) {
	mapQueue := pool.mapKeyQueue

	forList := pool.mapKeyQueue.List()

	for _, v := range forList {
		t, ok := pool.TCPList[v]

		if ok {
			interval := now.Sub(t.LastTime)
			isTimeout := false
			if interval > time.Second*5 {
				isTimeout = true
			}

			if isTimeout {
				pool.mapKeyQueue.RemoveValue(v)
				t.ResultData()
				delete(pool.TCPList, v)
			}
		} else {
			mapQueue.RemoveValue(v)
			log.Println("在时间队列中出现未知的连接Key ConversationPool.go 166")
		}
	}

}

func NewConnectionPool(dataChan chan<- []byte) *ConnectionPool {
	return &ConnectionPool{
		FragmentList: ip4defrag.NewIPv4Defragmenter(),
		connMsgs:     make(map[IPRefragKey]*ConnMsg),
		TCPList:      make(map[string]*TCPConnection),
		mapKeyQueue:  NewQueue(),
		DataChan:     dataChan,
	}
}
