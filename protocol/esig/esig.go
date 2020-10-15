package esig

import (
	"bytes"
	"errors"
	uuid "github.com/satori/go.uuid"
)

type ESIG [32]byte

func (e ESIG) High() uuid.UUID {

	return uuid.FromBytesOrNil(e[:16])
}

func (e ESIG) HaveHigh(high uuid.UUID) bool {

	return uuid.FromBytesOrNil(e[:16])
}

func (e ESIG) Low() uuid.UUID {

	return uuid.FromBytesOrNil(e[17:32])
}

// Equal returns true if u1 and u2 equals, otherwise returns false.
func Equal(u1 ESIG, u2 ESIG) bool {
	return bytes.Equal(u1[:], u2[:])
}

// Equal returns true if u1 and u2 equals, otherwise returns false.
func HighEqual(u1 ESIG, u2 ESIG) bool {
	return bytes.Equal(u1[:16], u2[:16])
}

var Nil = ESIG{}

func IsNul(e ESIG) bool {
	return Equal(e, Nil)
}

func FromUuid(u uuid.UUID) ESIG {

	var sig [32]byte
	copy(sig[:len(u)], u.Bytes())

	return sig
}

func From2Uuid(u1, u2 uuid.UUID) ESIG {

	var sig [32]byte
	copy(sig[:16], u1.Bytes())
	copy(sig[17:32], u2.Bytes())

	return sig
}

func FromBytes(b []byte) (ESIG, error) {

	size := len(b)
	valid := size == 16 || size == 32

	if !valid {
		return [32]byte{}, errors.New("wrong bytes len must be 16 or 32")
	}

	var sig [32]byte
	copy(sig[:size], b)

	return sig, nil
}