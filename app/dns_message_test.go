package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDeserializeLabels(t *testing.T) {
	labelsBytes := []byte{6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 100, 50, 100} //google.com label encoded(6google3com0) + random padding after null byte
	expected := "google.com"
	labels, bytesRead := deserializeLabels(labelsBytes, labelsBytes)
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
	labels, bytesRead := deserializeLabels(labelsBytes, labelsBytes)
	if bytesRead != 1 {
		t.Errorf("Expected 1 bytesRead but got %v", bytesRead)
	}
	if labels != expected {
		t.Errorf("Expected '%s' but got '%s'", expected, labels)
	}
}

func TestDeserializeLabelsWithPointer(t *testing.T) {
	labelsBytes := []byte{4, 103, 103, 103, 103, 6, 103, 111, 111, 103, 108, 101, 3, 99, 111, 109, 0, 100, 50, 100} //gggg.google.com label encoded(6google3com0) + random padding after null byte
	labelsBytesWithPointer := []byte{192, 5, 100, 100, 100}                                                         // pointer to "google.com" with offset 5 + padding

	expected := "google.com"
	labels, bytesRead := deserializeLabels(labelsBytesWithPointer, labelsBytes)
	if bytesRead != 2 {
		t.Errorf("Expected 2 bytesRead but got %v", bytesRead)
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
	dnsQuestions, _ := deserializeQuestions(in, 1, in)

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
	dnsQuestions, _ := deserializeQuestions(in, 2, in)

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

func TestDeserializeOneAnswer(t *testing.T) {
	in := []byte{12, 99, 111, 100, 101, 99, 114, 97, 102, 116, 101, 114, 115, 2, 105, 111, 0, 0, 1, 0, 1, 0, 0, 0, 60, 0, 4, 8, 8, 8, 8}
	expected := []DnsAnswer{{
		NAME:     "codecrafters.io",
		TYPE:     1,
		CLASS:    1,
		TTL:      60,
		RDLENGTH: 4,
		RDATA:    []byte{8, 8, 8, 8},
	}}

	dnsAnswers, _ := deserializeAnswers(in, 1, in)

	t.Logf("dnsAnswers: %v", dnsAnswers)
	if len(dnsAnswers) != 1 {
		t.Errorf("Expected only 1 answer but read %v", len(dnsAnswers))
	}
	if !reflect.DeepEqual(expected, dnsAnswers) {
		t.Errorf("Expected %v but got %v", expected, dnsAnswers)
	}
}

func TestDeserializeMultipleAnswers(t *testing.T) {
	in := []byte{12, 99, 111, 100, 101, 99, 114, 97, 102, 116, 101, 114, 115, 2, 105, 111, 0, 0, 1, 0, 1, 0, 0, 0, 60, 0, 4, 8, 8, 8, 8, 12, 99, 111, 100, 101, 99, 114,
		97, 102, 116, 101, 114, 115, 2, 105, 111, 0, 0, 1, 0, 1, 0, 0, 0, 120, 0, 4, 10, 10, 10, 10, 100, 100, 100, 100}
	expected := []DnsAnswer{{
		NAME:     "codecrafters.io",
		TYPE:     1,
		CLASS:    1,
		TTL:      60,
		RDLENGTH: 4,
		RDATA:    []byte{8, 8, 8, 8},
	}, {
		NAME:     "codecrafters.io",
		TYPE:     1,
		CLASS:    1,
		TTL:      120,
		RDLENGTH: 4,
		RDATA:    []byte{10, 10, 10, 10},
	}}

	dnsAnswers, offset := deserializeAnswers(in, 2, in)

	t.Logf("dnsAnswers: %v", dnsAnswers)
	if len(dnsAnswers) != 2 {
		t.Errorf("Expected 2 answers but read %v", len(dnsAnswers))
	}
	if !reflect.DeepEqual(expected, dnsAnswers) {
		t.Errorf("Expected %v but got %v", expected, dnsAnswers)
	}
	if offset != 62 {
		t.Errorf("Expected offset=66 but got %v", offset)
	}
}
