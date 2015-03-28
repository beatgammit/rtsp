package rtp

import (
	"net"
)

const (
	RTP_VERSION = 2
)

const (
	hasRtpPadding = 1 << 2
	hasRtpExt     = 1 << 3
)

// Packet as per https://tools.ietf.org/html/rfc1889#section-5.1
//
//  0                   1                   2                   3
//  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |V=2|P|X|  CC   |M|     PT      |       sequence number         |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |                           timestamp                           |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
// |           synchronization source (SSRC) identifier            |
// +=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+=+
// |            contributing source (CSRC) identifiers             |
// |                             ....                              |
// +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
type RtpPacket struct {
	Version        byte
	Padding        bool
	Ext            bool
	Marker         bool
	PayloadType    byte
	SequenceNumber uint
	Timestamp      uint
	SyncSource     uint

	CSRC []uint

	ExtHeader uint
	ExtData   []byte

	Payload []byte
}

type Session struct {
	Rtp  net.PacketConn
	Rtcp net.PacketConn

	RtpChan  <-chan RtpPacket
	RtcpChan <-chan []byte

	rtpChan  chan<- RtpPacket
	rtcpChan chan<- []byte
}

func New(rtp, rtcp net.PacketConn) *Session {
	rtpChan := make(chan RtpPacket, 10)
	rtcpChan := make(chan []byte, 10)
	s := &Session{
		Rtp:      rtp,
		Rtcp:     rtcp,
		RtpChan:  rtpChan,
		RtcpChan: rtcpChan,
		rtpChan:  rtpChan,
		rtcpChan: rtcpChan,
	}
	go s.HandleRtpConn(rtp)
	go s.HandleRtcpConn(rtcp)
	return s
}

func toUint(arr []byte) (ret uint) {
	for i, b := range arr {
		ret |= uint(b) << (8 * uint(len(arr)-i-1))
	}
	return ret
}

func (s *Session) HandleRtpConn(conn net.PacketConn) {
	buf := make([]byte, 4096)
	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}

		cpy := make([]byte, n)
		copy(cpy, buf)
		go s.handleRtp(cpy)
	}
}

func (s *Session) HandleRtcpConn(conn net.PacketConn) {
	buf := make([]byte, 4096)
	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			panic(err)
		}
		cpy := make([]byte, n)
		copy(cpy, buf)
		go s.handleRtcp(cpy)
	}
}

func (s *Session) handleRtp(buf []byte) {
	packet := RtpPacket{
		Version:        buf[0] & 0x03,
		Padding:        buf[0]&hasRtpPadding != 0,
		Ext:            buf[0]&hasRtpExt != 0,
		CSRC:           make([]uint, buf[0]>>4),
		Marker:         buf[1]&1 != 0,
		PayloadType:    buf[1] >> 1,
		SequenceNumber: toUint(buf[2:4]),
		Timestamp:      toUint(buf[4:8]),
		SyncSource:     toUint(buf[8:12]),
	}
	if packet.Version != RTP_VERSION {
		panic("Unsupported version")
	}

	i := 12

	for j := range packet.CSRC {
		packet.CSRC[j] = toUint(buf[i : i+4])
		i += 4
	}

	if packet.Ext {
		packet.ExtHeader = toUint(buf[i : i+2])
		length := toUint(buf[i+2 : i+4])
		i += 4
		if length > 0 {
			packet.ExtData = buf[i : i+int(length)*4]
			i += int(length) * 4
		}
	}

	packet.Payload = buf[i:]

	s.rtpChan <- packet
}

func (s *Session) handleRtcp(buf []byte) {
	// TODO: implement rtcp
}
