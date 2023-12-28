package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
)

func main() {
	resolver := parseArgs()

	udpConn, err := bindIp("127.0.0.1:2053")
	if err != nil {
		fmt.Println("error starting server. ", err)
	}
	defer udpConn.Close()

	var server Server
	if resolver != "" {
		server = NewForwardingServer(udpConn, resolver)
	} else {
		server = NewDefaultServer(udpConn)
	}
	err = server.Run()
	if err != nil {
		fmt.Println(err)
	}
}

func parseArgs() string {
	var resolver string
	flag.StringVar(&resolver, "resolver", "", "Optional DNS resolver to forward queries to. Should take the form <ip>:<port> or just <ip> (will use default dns port 53 in this case)")
	flag.Parse()

	if resolver == "" {
		return ""
	}

	ipPortSplit := strings.Split(resolver, ":")
	if len(ipPortSplit) == 1 {
		return ipPortSplit[0] + ":53"
	}
	return resolver
}

func bindIp(address string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return nil, err
	}
	return udpConn, nil
}
