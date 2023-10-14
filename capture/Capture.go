package capture

import (
	"errors"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

const (
	snapshotLen uint32 = 1526
	timeout            = 1 * time.Second
)

type Capture struct {
	Devices    []pcap.Interface
	handle     *pcap.Handle
	packets    chan gopacket.Packet
	sourceType uint
	EndChan    chan bool

	connectionPool *ConnectionPool
}

func NewCapture(dataChan chan<- []byte, endChan chan bool) (*Capture, error) {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return nil, err
	}
	return &Capture{
		Devices:        devices,
		connectionPool: NewConnectionPool(dataChan),
		EndChan:        endChan,
	}, nil
}

//打印网卡详细信息
func (capture Capture) PrintDevices() []string {
	var card []string

	fmt.Print("Devices found")
	for _, device := range capture.Devices {
		fmt.Println("\nName", device.Name)
		card = append(card, device.Name)
		fmt.Println("Description:", device.Description)
		for _, address := range device.Addresses {
			fmt.Println("- IP address:", address.IP)
			fmt.Println("- Subnet mask:", address.Netmask)
		}
	}
	return card
}

//sourceType:     1表示网卡， 0表示pcap文件
func (capture *Capture) SetCaptureSource(sourceName string,
	sourceType uint, promiscuous bool) error {

	capture.sourceType = sourceType
	//如果是嗅探网卡
	if sourceType == 1 {
		err := capture.setCaptureCard(sourceName, promiscuous)
		if err != nil {
			return err
		}

		packetSource := gopacket.NewPacketSource(capture.handle, capture.handle.LinkType())
		capture.packets = packetSource.Packets()
	} else if sourceType == 0 { //如果是分析文件
		err := capture.setCaptureFile(sourceName)
		if err != nil {
			return err
		}
		packetSource := gopacket.NewPacketSource(capture.handle, capture.handle.LinkType())
		capture.packets = packetSource.Packets()
	} else {
		return errors.New("sniffer源类型异常")
	}

	return nil
}

//设置嗅探网卡
func (capture *Capture) setCaptureCard(device string, promiscuous bool) error {
	handle, err := pcap.OpenLive(device, int32(snapshotLen), promiscuous, timeout)
	if err != nil {
		return err
	} else {
		capture.handle = handle
		return nil
	}

}

//分析文件
func (capture *Capture) setCaptureFile(fileName string) (err error) {
	capture.handle, err = pcap.OpenOffline(fileName)
	if err != nil {
		return
	}
	return
}

func (capture *Capture) StartCapture() {
	defer capture.handle.Close()

	if capture.sourceType == 1 {
		for {
			select {
			case packet := <-capture.packets:
				capture.connectionPool.DisposePacket(packet)
			case <-time.After(2 * time.Second):
				log.Println("2秒内没有连接到达...(sniffer.go 112)")
			case result := <-capture.EndChan:
				if result{
					for _, v := range capture.connectionPool.TCPList {
						v.ResultData()
					}
					return
				}
			}
		}
	} else {
		for packet := range capture.packets {
			capture.connectionPool.DisposePacket(packet)
		}

		//i:=1
		for _, v := range capture.connectionPool.TCPList {

			v.ResultData()
		}
		//println("end:",i)
		//close(capture.connectionPool.DataChan)
	}

}
