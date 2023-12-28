package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"net"
)

type ForwardingServer struct {
	conn     *net.UDPConn
	resolver string
}

func NewForwardingServer(udpConn *net.UDPConn, resolver string) *ForwardingServer {
	return &ForwardingServer{
		conn:     udpConn,
		resolver: resolver,
	}
}

func (server *ForwardingServer) Run() error {
	fmt.Println("Running forwarding server")

	resolverAddr, err := net.ResolveUDPAddr("udp", server.resolver)
	if err != nil {
		return fmt.Errorf("invalid resolver provided %v. Error is %v", server.resolver, err)
	}
	resolverConn, err := net.DialUDP("udp", nil, resolverAddr)
	if err != nil {
		return fmt.Errorf("error in UDP connection. Error is %v", err)
	}
	defer resolverConn.Close()

	clientBuf := make([]byte, 512)
	resolverBuf := make([]byte, 512)

	for {
		size, source, err := server.conn.ReadFromUDP(clientBuf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := clientBuf[:size]
		fmt.Printf("Request DNS message bytes: %08b\n", receivedData)
		fmt.Printf("Request DNS message hex from %s: %v\n", source, receivedData)

		receivedMessage := deserializeMessage(receivedData)
		fmt.Printf("Request DNS message: %+v\n", receivedMessage)

		responseMessage := parseMessage(receivedMessage, resolverConn, resolverBuf)
		fmt.Printf("Response DNS message: %+v\n", responseMessage)

		var response = responseMessage.serialize()
		fmt.Printf("Response DNS message bytes: %08b\n", response)
		fmt.Printf("Response DNS message hex: %v", response)

		_, err = server.conn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
	return nil
}

func parseMessage(receivedMessage *DnsMessage, resolverConn *net.UDPConn, resolverBuf []byte) *DnsMessage {
	rcvFlags := receivedMessage.Header.Flags
	rcvQuestions := receivedMessage.Questions

	resolverQuestionMsg, err := generateInitialQuestionMessage(rcvFlags.OPCODE, rcvFlags.RD)
	if err != nil {
		return generateResponseMessage(RCODE_SERVFAIL, receivedMessage)
	}

	responseMsg := generateResponseMessage(RCODE_NOERROR, receivedMessage)

	for _, question := range rcvQuestions {
		resolverQuestionMsg.Questions = []DnsQuestion{question}

		resolverConn.Write(resolverQuestionMsg.serialize())
		_, err := resolverConn.Read(resolverBuf)
		if err != nil {
			fmt.Println("error receiving data from resolver: ", err)
			return generateResponseMessage(RCODE_SERVFAIL, receivedMessage)
		}

		resolverRcvMessage := deserializeMessage(resolverBuf)
		resolverRcode := resolverRcvMessage.Header.Flags.RCODE
		if resolverRcode != RCODE_NOERROR {
			fmt.Printf("resolver DNS server returned RCODE %v, full DNS message is: %v\n", resolverRcode, resolverRcvMessage)
			return generateResponseMessage(resolverRcode, receivedMessage)
		}

		responseMsg.Header.ANCOUNT += resolverRcvMessage.Header.ANCOUNT
		responseMsg.Answers = append(responseMsg.Answers, resolverRcvMessage.Answers...)
	}
	return responseMsg
}

func generateInitialQuestionMessage(opcode uint8, rd bool) (*DnsMessage, error) {
	idBuf := make([]byte, 2)
	_, err := rand.Read(idBuf)
	if err != nil {
		return nil, err
	}
	return &DnsMessage{
		Header: DnsHeader{
			ID: binary.BigEndian.Uint16(idBuf),
			Flags: HeaderFlags{
				QR:     false,
				OPCODE: opcode,
				AA:     false,
				TC:     false,
				RD:     rd,
				RA:     false,
				Z:      0,
				RCODE:  0,
			},
			QDCOUNT: 1,
			ANCOUNT: 0,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
	}, nil
}

func generateResponseMessage(rcode uint8, receivedMessage *DnsMessage) *DnsMessage {
	return &DnsMessage{
		Header: DnsHeader{
			ID: receivedMessage.Header.ID,
			Flags: HeaderFlags{
				QR:     true,
				OPCODE: receivedMessage.Header.Flags.OPCODE,
				AA:     false,
				TC:     false,
				RD:     receivedMessage.Header.Flags.RD,
				RA:     false,
				Z:      0,
				RCODE:  rcode,
			},
			QDCOUNT: receivedMessage.Header.QDCOUNT,
			ANCOUNT: 0,
			NSCOUNT: 0,
			ARCOUNT: 0,
		},
		Questions: receivedMessage.Questions,
	}
}
