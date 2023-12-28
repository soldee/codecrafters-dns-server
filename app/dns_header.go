package main

import (
	"encoding/binary"
)

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
