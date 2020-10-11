package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/satori/go.uuid"
	"github.com/v8platform/rac/protocol/types"
	"github.com/xelaj/go-dry"
	"math"
)

const (
	UTF8_CHARSET   = "UTF-8"
	SIZEOF_SHORT   = 2
	SIZEOF_        = 4
	SIZEOF_LONG    = 8
	NULL_BYTE      = 0x80
	TRUE_BYTE      = 1
	FALSE_BYTE     = 0
	MAX_SHIFT      = 7
	NULL_SHIFT     = 6
	BYTE_MASK      = 255
	NEXT_MASK      = -128
	NULL_NEXT_MASK = 64
	LAST_MASK      = 0
	NULL_LSB_MASK  = 63
	LSB_MASK       = 127
	TEMP_CAPACITY  = 256
)

type encoder struct {
	*bytes.Buffer
}

func NewEncoder() *encoder {

	return &encoder{
		bytes.NewBuffer([]byte{}),
	}
}

func (e *encoder) encodeBoolean(val bool) {

	if val {
		_ = e.WriteByte(TRUE_BYTE)
	} else {
		_ = e.WriteByte(FALSE_BYTE)
	}

}

func (e *encoder) encodeByte(b byte) {

	_ = e.WriteByte(b)

}

func (e *encoder) encodeChar(val int) {

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(val))

	_, _ = e.Write(buf)

}

func (e *encoder) encodeShort(val int) {

	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(val))
	_, err := e.Write(buf)

	if err != nil {
		panic(err)
	}
}

func (e *encoder) encodeInt(val int) {

	e.encodeUint32(uint32(val))

}

func (e *encoder) encodeUint32(val uint32) {

	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	_, _ = e.Write(buf)
}

func (e *encoder) encodeUint64(val uint64) {

	buf := make([]byte, 4)
	binary.BigEndian.PutUint64(buf, val)
	_, _ = e.Write(buf)
}

func (e *encoder) encodeLong(val int64) {

	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(val))
	_, _ = e.Write(buf)
}

func (e *encoder) encodeFloat(val float32) {

	e.encodeUint32(math.Float32bits(val))

}

func (e *encoder) encodeDouble(val float64) {

	e.encodeUint64(math.Float64bits(val))

}

func (e *encoder) encodeNull() {

	e.encodeInt(0x80) // NULL_BYTE

}

func (e *encoder) encodeString(val string) {

	if len(val) == 0 {
		e.encodeNull()
		return
	}

	b := []byte(val)
	e.encodeNullableSize(len(b))
	e.Write(b)

}

func (e *encoder) encodeTypedValue(val interface{}) {

	if val == nil {
		e.encodeNull()
		return
	}

	ServiceWireType := detectType(val)
	e.encodeType(ServiceWireType)

	switch ServiceWireType {

	case BOOLEAN:
		e.encodeBoolean(val.(bool))
	case INT:
		e.encodeInt(val.(int))
	case LONG:
		e.encodeLong(val.(int64))
	case BYTE:
		e.encodeByte(val.(byte))
	default:
		dry.PanicIf(true, "error encode typed value")
	}

}

func detectType(val interface{}) ServiceWireType {

	switch val.(type) {
	case bool:
		return BOOLEAN
	case byte:
		return BYTE
	case int:
		return INT
	case int64:
		return LONG
	case uint64:
		return LONG
	default:
		return BYTE
	}

}

func (e *encoder) encodeUuid(val uuid.UUID) {

	b, _ := val.MarshalBinary()
	e.Write(b)

}

func (e *encoder) encodeSize(size int) {

	var b1 int

	msb := size >> MAX_SHIFT
	if msb != 0 {
		b1 = -128
	} else {
		b1 = 0
	}

	e.WriteByte(byte(b1 | (size & 0x7F)))

	for size = msb; size > 0; size = msb {

		msb >>= MAX_SHIFT
		if msb != 0 {
			b1 = -128
		} else {
			b1 = 0
		}

		e.WriteByte(byte(b1 | (size & 0x7F)))

	}

}

func (e *encoder) encodeNullableSize(size int) {

	var b1 int

	msb := size >> NULL_SHIFT
	if msb != 0 {
		b1 = 64
	} else {
		b1 = 0
	}

	_ = e.WriteByte(byte(b1 | (size & 0x3F)))

	for size = msb; size > 0; size = msb {

		msb >>= MAX_SHIFT
		if msb != 0 {
			b1 = -128
		} else {
			b1 = 0
		}

		_ = e.WriteByte(byte(b1 | (size & 0x7F)))

	}

}

func (e *encoder) encodeType(val types.Typed) {

	if int(val.Type()) == NULL_BYTE {
		e.encodeNull()
		return
	}
	e.WriteByte(byte(val.Type()))

}

func (e *encoder) encodeByteArray(val []byte) {

	e.Write(val)

}

func (e *encoder) encodeEndpointId(val int) {
	e.encodeNullableSize(val)
}

type Decoder struct {
	*bytes.Reader
}

func NewDecoder(buf []byte) *Decoder {

	return &Decoder{
		bytes.NewReader(buf),
	}
}

func (e *Decoder) decodeBoolean() bool {

	b, _ := e.ReadByte()

	switch b {

	case TRUE_BYTE:
		return true
	case FALSE_BYTE:
		return false
	default:
		panic(errors.New("error parse bool byte"))
	}

}

func (e *Decoder) decodeByte() byte {

	b, _ := e.ReadByte()
	return b
}

func (e *Decoder) decodeSlice(size int) [][]byte {

	var parts [][]byte
	part := make([]byte, size)

	for n, err := e.Read(part); n == size && err == nil; {

		parts = append(parts, part)

		//part = make([]byte, size)

	}

	return parts
}

func (e *Decoder) decodeUnsignedByte() byte {

	b, err := e.ReadByte()

	dry.PanicIfErr(err, e)

	return b
}

func (e *Decoder) decodeChar() uint16 {

	buf := make([]byte, 2)
	_, _ = e.Read(buf)
	char := binary.BigEndian.Uint16(buf)

	return char

}

func (e *Decoder) decodeShort() uint16 {

	buf := make([]byte, 2)

	n, _ := e.Read(buf)
	buf = buf[:n]

	char := binary.BigEndian.Uint16(buf)
	return char
}

func (e *Decoder) decodeUnsignedShort() uint16 {

	return e.decodeShort() & 0xFFFF

}

func (e *Decoder) decodeInt() int32 {

	return int32(e.decodeUint())
}

func (e *Decoder) decodeUint() uint32 {

	buf := make([]byte, 4)
	_, _ = e.Read(buf)
	char := binary.BigEndian.Uint32(buf)

	return char

}

func (e *Decoder) decodeLong() int64 {

	buf := make([]byte, 8)
	_, _ = e.Read(buf)
	char := binary.BigEndian.Uint64(buf)
	return int64(char)
}
func (e *Decoder) decodeFloat() float32 {

	b := e.decodeUint()
	return math.Float32frombits(b)
}

func (e *Decoder) decodeString() string {

	size := e.decodeNullableSize()
	if size == 0 {
		return ""
	}
	buf := make([]byte, size)
	_, err := e.Read(buf)

	dry.PanicIfErr(err, buf)

	return string(buf)

}

func (e *Decoder) decodeUUID() uuid.UUID {
	buf := make([]byte, 16)
	_, _ = e.Read(buf)
	u, _ := uuid.FromBytes(buf)
	return u
}

func (e *Decoder) decodeSize() int {

	ff := 0xFFFFFF80
	b1, _ := e.ReadByte()
	cur := int(b1 & 0xFF)
	size := cur & 0x7F
	for shift := MAX_SHIFT; (cur & ff) != 0x0; {

		b1, _ = e.ReadByte()
		cur = int(b1 & 0xFF)
		size += (cur & 0x7F) << shift
		shift += MAX_SHIFT
	}

	return size

}

func (e *Decoder) decodeNullableSize() int {
	size := 0
	//ff := 0xFFFFFF80
	b1, _ := e.ReadByte()
	cur := int(b1 & 0xFF)
	if (cur & 0xFFFFFF80) == 0x0 {
		size = cur & 0x3F
		if cur&0x40 == 0x0 {
			return size
		}

		shift := NULL_SHIFT
		b1, _ := e.ReadByte()
		cur := int(b1 & 0xFF)
		size += (cur & 0x7F) << NULL_SHIFT
		shift += MAX_SHIFT

		for (cur & 0xFFFFFF80) != 0x0 {

			b1, _ = e.ReadByte()
			cur = int(b1 & 0xFF)
			size += (cur & 0x7F) << shift
			shift += MAX_SHIFT

		}
		return size
	}

	if (cur & 0x7F) != 0x0 {
		panic("null expected")
	}

	return size

}

func (e *Decoder) decodeType() int {

	ff := 0xFFFFFF80
	b1, _ := e.ReadByte()

	cur := int(b1 & 0xFF)

	if cur&ff != 0x0 {

		if cur&0x7F != 0x0 {
			panic("null expected")
		}
		return int(NULL_TYPE)
	}

	return cur

}

func (e *Decoder) readUnsignedByte() uint16 {
	b, _ := e.ReadByte()
	return uint16(b & 0xFF)
}

func (e *Decoder) decodeEndpointId() int {

	id := e.decodeNullableSize()

	return id

}

func (e *Decoder) decodeByteArray() []byte {

	size := e.decodeSize()
	buf := make([]byte, size)

	_, _ = e.Read(buf)

	return buf

}
