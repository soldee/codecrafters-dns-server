package main

type DnsMessage struct {
	Header     DnsHeader
	Questions  []DnsQuestion
	Answers    []DnsAnswer
	Authority  DnsAuthority
	Additional DnsAdditional
}

type DnsAuthority struct {
}

type DnsAdditional struct {
}

func (msg *DnsMessage) serialize() []byte {
	msgBytes := append(
		msg.Header.serialize(),
		serializeQuestions(msg.Questions)...,
	)
	msgBytes = append(
		msgBytes,
		serializeAnswers(msg.Answers)...,
	)
	return msgBytes
}

func deserializeMessage(msgBytes []byte) *DnsMessage {
	dnsHeader := deserializeHeader(msgBytes[0:13])
	offsetAcc := 12
	dnsQuestions := deserializeQuestions(msgBytes, dnsHeader.QDCOUNT, &offsetAcc)
	dnsAnswers := deserializeAnswers(msgBytes, dnsHeader.ANCOUNT, &offsetAcc)

	return &DnsMessage{
		Header:    *dnsHeader,
		Questions: dnsQuestions,
		Answers:   dnsAnswers,
	}
}
