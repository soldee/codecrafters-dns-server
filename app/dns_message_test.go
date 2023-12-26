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
