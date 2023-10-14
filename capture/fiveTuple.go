package capture

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/google/gopacket/layers"
	"strconv"
)

const fnvBasis = 14695981039346656037
const fnvPrime = 1099511628211

type FiveTuple struct {
	SrcIP, DstIP     [4]byte
	SrcPort, DstPort uint16
	ProtocolType     layers.IPProtocol
}

func (f *FiveTuple) Hash() string {
	data := f.SrcIP[:]
	data = append(data, f.DstIP[:]...)
	data = append(data, uint16ToBytes(f.SrcPort)...)
	data = append(data, uint16ToBytes(f.DstPort)...)
	data = append(data, byte(f.ProtocolType))

	md5Inst :=  md5.New()
	md5Inst.Write(data)
	result := md5Inst.Sum([]byte(""))

	encodeStr := hex.EncodeToString(result)

	return encodeStr
}

func uint16ToBytes(num uint16) []byte {
	result := make([]byte, 0)
	result = append(result, byte(num>>8))
	result = append(result, byte(num&0xFF))
	return result
}

func IpToString(ip [4]byte) string {
	data := ""
	data += strconv.Itoa(int(ip[0])) + "."
	data += strconv.Itoa(int(ip[1])) + "."
	data += strconv.Itoa(int(ip[2])) + "."
	data += strconv.Itoa(int(ip[3]))
	return data
}
