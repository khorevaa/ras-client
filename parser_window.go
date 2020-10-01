package rac

import (
	"github.com/ericcornelissen/stringsx"
	"golang.org/x/text/encoding/charmap"
)

func decodeOutBytes(in []byte) ([]byte, error) {

	if stringsx.IsValidUTF8(string(in)) {
		return in, nil
	}

	dec := charmap.CodePage866.NewDecoder()
	out, err := dec.Bytes(in)

	return out, err
}
