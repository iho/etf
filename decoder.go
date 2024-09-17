package etf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

// DecodeErlTerm is the main entry point for decoding any ETF-encoded data.
func DecodeErlTerm(data []byte) (ErlTerm, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	if data[0] != ETFVersion {
		return nil, errors.New("invalid ETF version")
	}
	term, _, err := decodeTerm(data[1:])
	return term, err
}

// decodeTerm is the central decoding function that routes the decoding based on the type tag.
func decodeTerm(data []byte) (ErlTerm, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data to decode")
	}

	typeTag := data[0]
	switch typeTag {
	case ETF_ATOM:
		return decodeAtom(data[1:])
	case ETF_SMALL_INTEGER:
		return decodeSmallInteger(data[1:])
	case ETF_INTEGER:
		return decodeInteger(data[1:])
	case ETF_SMALL_TUPLE:
		return decodeSmallTuple(data[1:])
	case ETF_FLOAT:
		return decodeFloat(data[1:])
	case ETF_NIL:
		return Nil{}, 1, nil
	case ETF_LIST:
		return decodeList(data[1:])
	case ETF_BINARY:
		return decodeBinary(data[1:])
	case ETF_MAP:
		return decodeMap(data[1:])
	default:
		return nil, 0, fmt.Errorf("unsupported type tag: %d", typeTag)
	}
}

func decodeMap(data []byte) (ErlTerm, int, error) {
	if len(data) < 4 {
		return nil, 0, errors.New("insufficient data for map length")
	}
	length := int(binary.BigEndian.Uint32(data[:4]))
	terms := make(Map, 0, length)
	remaining := data[4:]
	totalBytes := 4
	for i := 0; i < length; i++ {
		key, keyBytes, err := decodeTerm(remaining)
		if err != nil {
			return nil, 0, err
		}
		value, valueBytes, err := decodeTerm(remaining[keyBytes:])
		if err != nil {
			return nil, 0, err
		}

		terms = append(terms, MapElem{Key: key, Value: value}) // store key-value pair in terms map
		remaining = remaining[keyBytes+valueBytes:]
		totalBytes += keyBytes + valueBytes
	}
	return terms, totalBytes, nil // return totalBytes without +1
}

// decodeAtom decodes an atom from the data.
func decodeAtom(data []byte) (ErlTerm, int, error) {
	if len(data) < 1 {
		return nil, 0, errors.New("insufficient data for atom length")
	}
	length := int(data[0])
	if len(data[1:]) < length {
		return nil, 0, errors.New("insufficient data for atom value")
	}
	return Atom(string(data[1 : 1+length])), 2 + length, nil
}

// decodeSmallInteger decodes a small integer.
func decodeSmallInteger(data []byte) (ErlTerm, int, error) {
	if len(data) < 1 {
		return nil, 0, errors.New("insufficient data for small integer")
	}
	return Integer(data[0]), 2, nil
}

// decodeInteger decodes a 32-bit integer.
func decodeInteger(data []byte) (ErlTerm, int, error) {
	if len(data) < 4 {
		return nil, 0, errors.New("insufficient data for integer")
	}
	value := int32(binary.BigEndian.Uint32(data[:4]))
	return Integer(value), 5, nil
}

// decodeSmallTuple decodes a small tuple and recursively decodes its elements.
func decodeSmallTuple(data []byte) (ErlTerm, int, error) {
	if len(data) < 1 {
		return nil, 0, errors.New("insufficient data for small tuple arity")
	}
	arity := int(data[0])
	terms := make(Tuple, arity)

	remaining := data[1:]
	totalBytes := 1
	for i := 0; i < arity; i++ {
		term, bytesRead, err := decodeTerm(remaining)
		if err != nil {
			return nil, 0, err
		}
		terms[i] = term
		remaining = remaining[bytesRead:]
		totalBytes += bytesRead
	}
	return Tuple(terms), totalBytes + 1, nil
}

// decodeFloat decodes a float.
func decodeFloat(data []byte) (ErlTerm, int, error) {
	if len(data) < 8 {
		return nil, 0, errors.New("insufficient data for float")
	}
	value := binary.BigEndian.Uint64(data[:8])
	return Float(math.Float64frombits(value)), 9, nil
}

// decodeList decodes a list and recursively decodes its elements.
func decodeList(data []byte) (ErlTerm, int, error) {
	if len(data) < 4 {
		return nil, 0, errors.New("insufficient data for list length")
	}
	length := int(binary.BigEndian.Uint32(data[:4]))
	terms := make(List, length)
	remaining := data[4:]
	totalBytes := 4
	for i := 0; i < length; i++ {
		term, bytesRead, err := decodeTerm(remaining)
		if err != nil {
			return nil, 0, err
		}
		terms[i] = term
		remaining = remaining[bytesRead:]
		totalBytes += bytesRead
	}
	return terms, totalBytes + 1, nil
}

// decodeBinary decodes a binary and returns the data as a Binary term.
func decodeBinary(data []byte) (ErlTerm, int, error) {
	if len(data) < 4 {
		return nil, 0, errors.New("insufficient data for binary length")
	}
	length := int(binary.BigEndian.Uint32(data[:4]))
	if len(data[4:]) < length {
		return nil, 0, errors.New("insufficient data for binary")
	}
	return Binary(data[4 : 4+length]), 5 + length, nil
}
