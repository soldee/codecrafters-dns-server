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
}

type DnsAuthority struct {
}

type DnsAdditional struct {
}

func (msg *DnsMessage) serialize() []byte {
	return append(
		msg.Header.serialize(),
		msg.Question.serialize()...,
	)
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

	questionBytes = append(questionBytes, question.serializeLabels()...)
	binary.BigEndian.AppendUint16(questionBytes, question.QTYPE)
	binary.BigEndian.AppendUint16(questionBytes, question.QCLASS)

	return questionBytes
}

func (question *DnsQuestion) serializeLabels() []byte {
	var labelsBytes []byte

	labels := strings.Split(question.QNAME, ".")
	for _, label := range labels {
		labelsBytes = append(labelsBytes, uint8(len(label)))
		labelsBytes = append(labelsBytes, label...)
	}
	labelsBytes = append(labelsBytes, 0x00)
	return labelsBytes
}
