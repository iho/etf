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
	default:
		return nil, errors.New("unsupported Erlang term type")
	}
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
	fmt.Println("tuple", tuple)
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
		fmt.Println("term", term)
		termBytes, err := EncodeErlTerm(term, false)
		if err != nil {
			return nil, err
		}
		fmt.Println("termBytes", termBytes)
		buf.Write(termBytes)
	}
	return buf.Bytes(), nil
}