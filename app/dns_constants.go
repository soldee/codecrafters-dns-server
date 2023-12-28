package main

const (
	TYPE_A     uint16 = 1
	TYPE_NS    uint16 = 2
	TYPE_MD    uint16 = 3
	TYPE_MF    uint16 = 4
	TYPE_CNAME uint16 = 5
	TYPE_SOA   uint16 = 6
	TYPE_MB    uint16 = 7
	TYPE_MG    uint16 = 8
	TYPE_MR    uint16 = 9
	TYPE_NULL  uint16 = 10
	TYPE_WKS   uint16 = 11
	TYPE_PTR   uint16 = 12
	TYPE_HINFO uint16 = 13
	TYPE_MINFO uint16 = 14
	TYPE_MX    uint16 = 15
	TYPE_TXT   uint16 = 16
)

const (
	CLASS_IN = 1
	CLASS_CS = 2
	CLASS_CH = 3
	CLASS_HS = 4
)

const (
	RCODE_NOERROR  = 0 // DNS Query completed successfully
	RCODE_FORMERR  = 1 // DNS Query Format Error
	RCODE_SERVFAIL = 2 // Server failed to complete the DNS request
	RCODE_NXDOMAIN = 3 // Domain name does not exist.
	RCODE_NOTIMP   = 4 // Function not implemented
	RCODE_REFUSED  = 5 // The server refused to answer for the query
	RCODE_YXDOMAIN = 6 // Name that should not exist, does exist
	RCODE_XRRSET   = 7 // RRset that should not exist, does exist
	RCODE_NOTAUTH  = 8 // Server not authoritative for the zone
	RCODE_NOTZONE  = 9 // Name not in zone
)
