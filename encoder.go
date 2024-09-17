package etf

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// Constants for ETF version and term types
const (
	ETFVersion        byte = 131
	ETF_ATOM          byte = 119
	ETF_SMALL_INTEGER byte = 97
	ETF_INTEGER       byte = 98
	ETF_FLOAT         byte = 70
	ETF_SMALL_TUPLE   byte = 104
	ETF_LARGE_TUPLE   byte = 105
	ETF_NIL           byte = 106
	ETF_LIST          byte = 108
	ETF_BINARY        byte = 109
	ETF_MAP           byte = 116
)

func EncodeErlTerm(term ErlTerm, writeHeader bool) ([]byte, error) {
	buf := new(bytes.Buffer)
	if writeHeader {
		buf.WriteByte(ETFVersion)
	}
	switch v := term.(type) {
	case Atom:
		return encodeAtom(buf, v)
	case Integer:
		return encodeInteger(buf, v)
	case Float:
		return encodeFloat(buf, v)
	case Tuple:
		return encodeTuple(buf, v)
	case List:
		return encodeList(buf, v)
	case Binary:
		return encodeBinary(buf, v)
	case Map:
		return encodeMap(buf, v)
	case Nil:
		return writeNil(buf)
	default:
		fmt.Println("Unsupported type:", v)
		return nil, errors.New("unsupported Erlang term type")
	}
}

func encodeMap(buf *bytes.Buffer, m Map) ([]byte, error) {
	buf.WriteByte(ETF_MAP)
	// Encode as 4-byte big endian
	length := uint32(len(m))
	err := binary.Write(buf, binary.BigEndian, length)
	if err != nil {
		return nil, err
	}
	for _, v := range m {
		_, ok := v.Key.(Nil)
		if ok {
			return nil, errors.New("map key cannot be nil")
		}
		keyBytes, err := EncodeErlTerm(v.Key, false)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		termBytes, err := EncodeErlTerm(v.Value, false)
		if err != nil {
			return nil, err
		}
		buf.Write(termBytes)
	}
	return buf.Bytes(), nil
}

func writeNil(buf *bytes.Buffer) ([]byte, error) {
	buf.WriteByte(ETF_NIL)
	return buf.Bytes(), nil
}

func encodeBinary(buf *bytes.Buffer, bin Binary) ([]byte, error) {
	buf.WriteByte(ETF_BINARY)
	// Encode as 4-byte big endian
	length := uint32(len(bin))
	err := binary.Write(buf, binary.BigEndian, length)
	if err != nil {
		return nil, err
	}
	buf.Write(bin)
	return buf.Bytes(), nil
}

func encodeList(buf *bytes.Buffer, list List) ([]byte, error) {
	buf.WriteByte(ETF_LIST)
	// Encode as 4-byte big endian
	length := uint32(len(list))
	err := binary.Write(buf, binary.BigEndian, length)
	if err != nil {
		return nil, err
	}
	for _, term := range list {
		termBytes, err := EncodeErlTerm(term, false)
		if err != nil {
			return nil, err
		}
		buf.Write(termBytes)
	}
	return buf.Bytes(), nil
}

func encodeAtom(buf *bytes.Buffer, atom Atom) ([]byte, error) {
	buf.WriteByte(ETF_ATOM)
	// ETF atom format: 2 bytes length + atom bytes
	atomBytes := []byte(atom)
	if len(atomBytes) > 65535 {
		return nil, errors.New("atom too long")
	}
	length := uint8(len(atomBytes))
	err := binary.Write(buf, binary.BigEndian, length)
	if err != nil {
		return nil, err
	}
	buf.Write(atomBytes)

	return buf.Bytes(), nil
}

func encodeInteger(buf *bytes.Buffer, integer Integer) ([]byte, error) {
	if integer >= 0 && integer <= 255 {
		buf.WriteByte(ETF_SMALL_INTEGER)
		buf.WriteByte(byte(integer))
	} else {
		buf.WriteByte(ETF_INTEGER)
		// Encode as 4-byte big endian
		err := binary.Write(buf, binary.BigEndian, int32(integer))
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func encodeFloat(buf *bytes.Buffer, floatVal Float) ([]byte, error) {
	buf.WriteByte(ETF_FLOAT)
	if err := binary.Write(buf, binary.BigEndian, floatVal); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodeTuple(buf *bytes.Buffer, tuple Tuple) ([]byte, error) {
	if len(tuple) <= 255 {
		buf.WriteByte(ETF_SMALL_TUPLE)
		buf.WriteByte(byte(len(tuple)))
	} else {
		buf.WriteByte(ETF_LARGE_TUPLE)
		// Encode as 4-byte big endian
		length := uint32(len(tuple))
		err := binary.Write(buf, binary.BigEndian, length)
		if err != nil {
			return nil, err
		}
	}
	for _, term := range tuple {
		termBytes, err := EncodeErlTerm(term, false)
		if err != nil {
			return nil, err
		}
		buf.Write(termBytes)
	}
	return buf.Bytes(), nil
}
