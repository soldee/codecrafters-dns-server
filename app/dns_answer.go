package main

import (
	"encoding/binary"
)

type DnsAnswer struct {
	NAME     string
	TYPE     uint16
	CLASS    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

func serializeAnswers(answers []DnsAnswer) []byte {
	var answerBytes []byte
	for _, answer := range answers {
		answerBytes = append(answerBytes, serializeLabels(answer.NAME)...)
		answerBytes = binary.BigEndian.AppendUint16(answerBytes, answer.TYPE)
		answerBytes = binary.BigEndian.AppendUint16(answerBytes, answer.CLASS)
		answerBytes = binary.BigEndian.AppendUint32(answerBytes, answer.TTL)
		answerBytes = binary.BigEndian.AppendUint16(answerBytes, answer.RDLENGTH)
		answerBytes = append(answerBytes, answer.RDATA...)
	}
	return answerBytes
}

func deserializeAnswers(msgBytes []byte, anCount uint16, offset *int) []DnsAnswer {
	answers := make([]DnsAnswer, anCount)

	for i := 0; i < int(anCount); i++ {
		labels, labelsOffset := deserializeLabels(msgBytes[*offset:], msgBytes[:*offset])
		*offset += labelsOffset
		rdlength := binary.BigEndian.Uint16(msgBytes[*offset+8 : *offset+10])
		answers[i] = DnsAnswer{
			NAME:     labels,
			TYPE:     binary.BigEndian.Uint16(msgBytes[*offset : *offset+2]),
			CLASS:    binary.BigEndian.Uint16(msgBytes[*offset+2 : *offset+4]),
			TTL:      binary.BigEndian.Uint32(msgBytes[*offset+4 : *offset+8]),
			RDLENGTH: rdlength,
			RDATA:    msgBytes[*offset+10 : *offset+10+int(rdlength)],
		}
		*offset += 10 + int(rdlength)
	}
	return answers
}
