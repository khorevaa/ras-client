package ras_client

import "bytes"

func NewBuffer() *bytes.Buffer {
	return bytes.NewBuffer([]byte{})
}

func wrapBuffer(buf ...*bytes.Buffer) *bytes.Buffer {

	header := NewBuffer()

	for _, buffer := range buf {
		buffer.WriteTo(header)
	}

	return header

}

func detectSupportedVersion(err error) string {

	// todo
	return "9.0"

}
