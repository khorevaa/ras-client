package serialize

import (
	uuid "github.com/satori/go.uuid"
	"io"
	"time"
)

type BinaryWriter interface {
	Format(codec Encoder, version int, writer io.Writer)
}

type BinaryParser interface {
	Parse(codec Decoder, version int, reader io.Reader)
}

type Encoder interface {
	Bool(val bool, w io.Writer)
	Byte(val byte, w io.Writer)
	Char(val int, w io.Writer)
	Short(val int16, w io.Writer)
	Int(val int, w io.Writer)
	Uint(val uint, w io.Writer)
	Int16(val int16, w io.Writer)
	Uint16(val uint16, w io.Writer)
	Int32(val int32, w io.Writer)
	Uint32(val uint32, w io.Writer)
	Int64(val int64, w io.Writer)
	Uint64(val uint64, w io.Writer)
	// Long is copy og Int64
	Long(val int64, w io.Writer)
	Float32(val float32, w io.Writer)
	Float64(val float64, w io.Writer)
	// Double is copy og Float64
	Double(val float64, w io.Writer)

	Null(w io.Writer)
	String(val string, w io.Writer)
	TypedValue(val interface{}, w io.Writer)
	Uuid(val uuid.UUID, w io.Writer)
	Size(val int, w io.Writer)
	NullableSize(val int, w io.Writer)
	Type(val byte, w io.Writer)
	EndpointId(val int, w io.Writer)
	Time(val time.Time, w io.Writer)
	// Bytes is alias ByteArray
	Bytes(val []byte, w io.Writer)
	Value(val interface{}, w io.Writer)
}

type Decoder interface {
	BoolPtr(val *bool, r io.Reader)
	Bool(r io.Reader) (bool, bool)

	BytePtr(val *byte, r io.Reader)
	Byte(r io.Reader) byte

	CharPtr(ptr *int16, r io.Reader)
	Char(r io.Reader) int16

	ShortPtr(val *int16, r io.Reader)
	Short(r io.Reader) int16

	IntPtr(val *int, r io.Reader)
	Int(r io.Reader) int

	UintPtr(val *uint, r io.Reader)
	Uint(r io.Reader) uint

	Uint16(r io.Reader) uint16
	Uint16Ptr(ptr *uint16, r io.Reader)

	Int32Ptr(val *int32, r io.Reader)
	Int32(r io.Reader) int32

	Uint32Ptr(val *uint32, r io.Reader)
	Uint32(r io.Reader) uint32

	Int64Ptr(val *int64, r io.Reader)
	Int64(r io.Reader) int64

	Uint64Ptr(val *uint64, r io.Reader)
	Uint64(r io.Reader) uint64

	// Long is copy og Int64
	LongPtr(val *int64, r io.Reader)
	Long(r io.Reader) int64

	Float32Ptr(val *float32, r io.Reader)
	Float32(r io.Reader) float32

	Float64Ptr(val *float64, r io.Reader)
	Float64(r io.Reader) float64
	// Double is copy og Float64
	DoublePtr(val *float64, r io.Reader)
	Double(r io.Reader) float64

	Null(r io.Reader)

	StringPtr(val *string, r io.Reader)
	String(r io.Reader) string

	TypedValue(val interface{}, r io.Reader)

	UuidPtr(val *uuid.UUID, r io.Reader)
	Uuid(r io.Reader) uuid.UUID

	Size(r io.Reader) int
	NullableSize(r io.Reader) int
	Type(r io.Reader) byte
	// Bytes is alias ByteArray
	Bytes(val []byte, r io.Reader)

	EndpointIdPtr(ptr *int, r io.Reader)
	EndpointId(r io.Reader) int

	TimePtr(ptr *time.Time, r io.Reader)
	Time(r io.Reader) time.Time

	Value(val interface{}, r io.Reader)
}
