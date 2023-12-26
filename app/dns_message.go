package main

import (
	"encoding/binary"
	"strings"
)

type DnsMessage struct {
	Header     DnsHeader
	Question   DnsQuestion
	Answer     DnsAnswer
	Authority  DnsAuthority
	Additional DnsAdditional
}

type DnsHeader struct {
	ID      uint16
	Flags   HeaderFlags
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

type HeaderFlags struct {
	QR     bool
	OPCODE uint8 // 4 bits
	AA     bool
	TC     bool
	RD     bool
	RA     bool
	Z      uint8 // 3 bits
	RCODE  uint8 // 4 bits
}

type DnsQuestion struct {
	QNAME  string
	QTYPE  uint16
	QCLASS uint16
}

type DnsAnswer struct {
	NAME     string
	TYPE     uint16
	CLASS    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

type DnsAuthority struct {
}

type DnsAdditional struct {
}

func (msg *DnsMessage) serialize() []byte {
	msgBytes := append(
		msg.Header.serialize(),
		msg.Question.serialize()...,
	)
	msgBytes = append(
		msgBytes,
		msg.Answer.serialize()...,
	)
	return msgBytes
}

func (header *DnsHeader) serialize() []byte {
	headerBytes := make([]byte, 12)

	binary.BigEndian.PutUint16(headerBytes, header.ID)
	binary.BigEndian.PutUint16(headerBytes[2:4], header.Flags.serialize())
	binary.BigEndian.PutUint16(headerBytes[4:6], header.QDCOUNT)
	binary.BigEndian.PutUint16(headerBytes[6:8], header.ANCOUNT)
	binary.BigEndian.PutUint16(headerBytes[8:10], header.NSCOUNT)
	binary.BigEndian.PutUint16(headerBytes[10:12], header.ARCOUNT)

	return headerBytes
}

func (flags *HeaderFlags) serialize() uint16 {
	var flagsBytes uint16

	if flags.QR {
		flagsBytes++
	}

	flagsBytes = flagsBytes << 4
	flagsBytes = flagsBytes | uint16(flags.OPCODE&0x07)

	flagsBytes = flagsBytes << 1
	if flags.AA {
		flagsBytes++
	}
	flagsBytes = flagsBytes << 1
	if flags.TC {
		flagsBytes++
	}
	flagsBytes = flagsBytes << 1
	if flags.RD {
		flagsBytes++
	}
	flagsBytes = flagsBytes << 1
	if flags.RA {
		flagsBytes++
	}

	flagsBytes = flagsBytes << 3
	flagsBytes = flagsBytes | uint16(flags.Z&0x07)

	flagsBytes = flagsBytes << 4
	flagsBytes = flagsBytes | uint16(flags.RCODE&0x0F)

	return flagsBytes
}

func (question *DnsQuestion) serialize() []byte {
	var questionBytes []byte

	questionBytes = append(questionBytes, serializeLabels(question.QNAME)...)
	questionBytes = binary.BigEndian.AppendUint16(questionBytes, question.QTYPE)
	questionBytes = binary.BigEndian.AppendUint16(questionBytes, question.QCLASS)

	return questionBytes
}

func serializeLabels(labelsStr string) []byte {
	var labelsBytes []byte

	labels := strings.Split(labelsStr, ".")
	for _, label := range labels {
		labelsBytes = append(labelsBytes, uint8(len(label)))
		labelsBytes = append(labelsBytes, label...)
	}
	labelsBytes = append(labelsBytes, 0x00)
	return labelsBytes
}

func (answer *DnsAnswer) serialize() []byte {
	var answerBytes []byte

	answerBytes = append(answerBytes, serializeLabels(answer.NAME)...)
	answerBytes = binary.BigEndian.AppendUint16(answerBytes, answer.TYPE)
	answerBytes = binary.BigEndian.AppendUint16(answerBytes, answer.CLASS)
	answerBytes = binary.BigEndian.AppendUint32(answerBytes, answer.TTL)
	answerBytes = binary.BigEndian.AppendUint16(answerBytes, answer.RDLENGTH)
	answerBytes = append(answerBytes, answer.RDATA...)
	return answerBytes
}

func deserializeMessage(msgBytes []byte) *DnsMessage {
	dnsHeader := deserializeHeader(msgBytes[0:13])
	offsetAcc := 12
	dnsQuestion, questionOffset := deserializeQuestion(msgBytes[offsetAcc:])
	offsetAcc += questionOffset
	dnsAnswer, _ := deserializeAnswer(msgBytes[offsetAcc:])

	return &DnsMessage{
		Header:   *dnsHeader,
		Question: *dnsQuestion,
		Answer:   *dnsAnswer,
	}
}

func deserializeHeader(headerBytes []byte) *DnsHeader {
	return &DnsHeader{
		ID:      binary.BigEndian.Uint16(headerBytes[0:2]),
		Flags:   *deserializeHeaderFlags(headerBytes[2:4]),
		QDCOUNT: binary.BigEndian.Uint16(headerBytes[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(headerBytes[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(headerBytes[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(headerBytes[10:12]),
	}
}

func deserializeHeaderFlags(flagsBytes []byte) *HeaderFlags {
	return &HeaderFlags{
		QR:     (flagsBytes[0] & 0x80) != 0,
		OPCODE: flagsBytes[0] & 0x78 >> 3,
		AA:     (flagsBytes[0] & 0x4) != 0,
		TC:     (flagsBytes[0] & 0x2) != 0,
		RD:     (flagsBytes[0] & 0x1) != 0,
		RA:     (flagsBytes[1] & 0x80) != 0,
		Z:      flagsBytes[1] & 0x70 >> 4,
		RCODE:  flagsBytes[1] & 0x0F,
	}
}

func deserializeQuestion(questionBytes []byte) (*DnsQuestion, int) {
	labels, labelsOffset := deserializeLabels(questionBytes)
	return &DnsQuestion{
		QNAME:  labels,
		QTYPE:  binary.BigEndian.Uint16(questionBytes[labelsOffset : labelsOffset+2]),
		QCLASS: binary.BigEndian.Uint16(questionBytes[labelsOffset+2 : labelsOffset+4]),
	}, labelsOffset + 4
}

func deserializeLabels(labelsBytes []byte) (string, int) {
	//6google3com0
	var labels string
	var offset int = 1
	var labelBytesLeft uint8 = 0
	for i, b := range labelsBytes {
		if b == 0x00 {
			break
		}
		offset++
		if labelBytesLeft == 0 {
			labelBytesLeft = uint8(b)
			if i != 0 {
				labels += "."
			}
		} else {
			labels += string(b)
			labelBytesLeft--
		}
	}
	return labels, offset
}

func deserializeAnswer(answerBytes []byte) (*DnsAnswer, int) {
	labels, offset := deserializeLabels(answerBytes)
	rdlength := binary.BigEndian.Uint16(answerBytes[offset+8 : offset+10])
	return &DnsAnswer{
		NAME:     labels,
		TYPE:     binary.BigEndian.Uint16(answerBytes[offset : offset+2]),
		CLASS:    binary.BigEndian.Uint16(answerBytes[offset+2 : offset+4]),
		TTL:      binary.BigEndian.Uint32(answerBytes[offset+4 : offset+8]),
		RDLENGTH: rdlength,
		RDATA:    answerBytes[offset+10 : offset+10+int(rdlength)],
	}, offset + 10 + int(rdlength) + 1
}
