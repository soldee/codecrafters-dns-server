package main

import (
	"encoding/binary"
)

type DnsQuestion struct {
	QNAME  string
	QTYPE  uint16
	QCLASS uint16
}

func serializeQuestions(questions []DnsQuestion) []byte {
	var questionsBytes []byte
	for _, question := range questions {
		questionsBytes = append(questionsBytes, serializeLabels(question.QNAME)...)
		questionsBytes = binary.BigEndian.AppendUint16(questionsBytes, question.QTYPE)
		questionsBytes = binary.BigEndian.AppendUint16(questionsBytes, question.QCLASS)
	}
	return questionsBytes
}

func deserializeQuestions(msgBytes []byte, qdCount uint16, offset *int) []DnsQuestion {
	questions := make([]DnsQuestion, qdCount)

	for i := 0; i < int(qdCount); i++ {
		labels, labelsOffset := deserializeLabels(msgBytes[*offset:], msgBytes[:*offset])
		*offset += labelsOffset
		questions[i] = DnsQuestion{
			QNAME:  labels,
			QTYPE:  binary.BigEndian.Uint16(msgBytes[*offset : *offset+2]),
			QCLASS: binary.BigEndian.Uint16(msgBytes[*offset+2 : *offset+4]),
		}
		*offset += 4
	}
	return questions
}
