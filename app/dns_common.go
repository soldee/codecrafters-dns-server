package main

import (
	"strings"
)

func serializeLabels(labelsStr string) []byte {
	var labelsBytes []byte

	labels := strings.Split(labelsStr, ".")
	for _, label := range labels {
		labelsBytes = append(labelsBytes, uint8(len(label)))
		labelsBytes = append(labelsBytes, label...)
	}
	labelsBytes = append(labelsBytes, 0x00)
	return labelsBytes
}

func deserializeLabels(labelsBytes []byte, readBytes []byte) (string, int) {
	var labels string
	var offset int = 1
	var labelBytesLeft uint8 = 0
	var pointer uint16 = 0
	for i, b := range labelsBytes {
		if b == 0x00 {
			break
		}
		offset++
		if labelBytesLeft == 0 {
			labelBytesLeft = uint8(b)
			if i != 0 {
				labels += "."
			}
			if labelBytesLeft&0xC0 == 0xC0 {
				pointer = (uint16(labelBytesLeft&0x3F) << 8) | uint16(labelsBytes[offset-1])
				pointerLabel, _ := deserializeLabels(readBytes[pointer:], readBytes)
				labels += pointerLabel
				break
			}
		} else {
			labels += string(b)
			labelBytesLeft--
		}
	}
	return labels, offset
}
