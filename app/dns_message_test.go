package main

import (
	"bytes"
	"testing"
)

func TestDeserializeLabels(t *testing.T) {
	labelsBytes := []byte{6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 100, 50, 100} //google.com label encoded(6google3com0) + random padding after null byte
	expected := "google.com"
	labels, bytesRead := deserializeLabels(labelsBytes)
	if bytesRead != 12 {
		t.Errorf("Expected 12 bytesRead but got %v", bytesRead)
	}
	if labels != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, labels)
	}
}

func TestDeserializeEmptyLabels(t *testing.T) {
	labelsBytes := []byte{0} // null byte only
	expected := ""
	labels, bytesRead := deserializeLabels(labelsBytes)
	if bytesRead != 1 {
		t.Errorf("Expected 1 bytesRead but got %v", bytesRead)
	}
	if labels != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, labels)
	}
}

func TestSerializeLabels(t *testing.T) {
	expectedBytes := []byte{6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0} //google.com label encoded
	labelsBytes := serializeLabels("google.com")
	if !bytes.Equal(labelsBytes, []byte(expectedBytes)) {
		t.Errorf("Expected '%v' but got '%v'", expectedBytes, labelsBytes)
	}
}

func TestDeserializeOneQuestion(t *testing.T) {
	in := []byte{6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1}
	expected := DnsQuestion{
		QNAME:  "google.com",
		QTYPE:  TYPE_A,
		QCLASS: CLASS_IN,
	}
	dnsQuestions, _ := deserializeQuestions(in, 1)

	if len(dnsQuestions) != 1 {
		t.Errorf("Expected only 1 question but read %v", len(dnsQuestions))
		t.Logf("Questions deserialized: %v", dnsQuestions)
	}
	if dnsQuestions[0] != expected {
		t.Errorf("Expected %v but got %v", expected, dnsQuestions[0])
	}
}

func TestDeserializeMultipleQuestions(t *testing.T) {
	in := []byte{6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 0, 1, 0, 1, 12, 99, 111, 100, 101, 99, 114, 97, 102, 116, 101, 114, 115, 2, 105, 111, 0, 0, 5, 0, 4}
	expected := []DnsQuestion{{
		QNAME:  "google.com",
		QTYPE:  TYPE_A,
		QCLASS: CLASS_IN,
	}, {
		QNAME:  "codecrafters.io",
		QTYPE:  TYPE_CNAME,
		QCLASS: CLASS_HS,
	}}
	dnsQuestions, _ := deserializeQuestions(in, 2)

	if len(dnsQuestions) != 2 {
		t.Errorf("Expected 2 questions but read %v", len(dnsQuestions))
		t.Logf("Questions deserialized: %v", dnsQuestions)
	}
	if dnsQuestions[0] != expected[0] {
		t.Errorf("Expected %v but got %v", expected, dnsQuestions)
	}
	if dnsQuestions[1] != expected[1] {
		t.Errorf("Expected %v but got %v", expected, dnsQuestions)
	}
}
