package etf

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

func DecodeErlTerm(data []byte) (ErlTerm, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}
	fmt.Printf("data[0]: %v \n", data[0])
	if data[0] != ETFVersion {
		return nil, errors.New("invalid ETF version")
	}

	if len(data) < 2 {
		return nil, errors.New("insufficient data for type tag")
	}

	typeTag := data[1]

	switch typeTag {
	case ETF_ATOM:
		return decodeAtom(data[2:])
	case ETF_SMALL_INTEGER:
		return decodeSmallInteger(data[2:])
	case ETF_INTEGER:
		return decodeInteger(data[2:])
	case ETF_SMALL_TUPLE:
		return decodeSmallTuple(data[2:])
	case ETF_FLOAT:
		return decodeFloat(data[2:])
	default:
		return nil, fmt.Errorf("unsupported type tag: %d", typeTag)
	}
}

func decodeAtom(data []byte) (ErlTerm, error) {
	if len(data) < 2 {
		return nil, errors.New("insufficient data for atom length")
	}
	fmt.Println("data: ", string(data))
	fmt.Println("data: ", data)
	length := data[0]

	fmt.Println("length: ", length)
	// if len(data[1:]) != int(length) {
	// 	return nil, errors.New("insufficient data for atom")
	// }
	fmt.Println("data[2:]: ", string(data[1:1+length]))
	atom := Atom(string(data[1 : 1+length]))
	return atom, nil
}

func decodeSmallInteger(data []byte) (ErlTerm, error) {
	if len(data) < 1 {
		return nil, errors.New("insufficient data for small integer")
	}
	value := Integer(data[0])
	return value, nil
}

func decodeInteger(data []byte) (ErlTerm, error) {
	if len(data) < 4 {
		return nil, errors.New("insufficient data for integer")
	}
	value := int32(binary.BigEndian.Uint32(data[:4]))
	return Integer(value), nil
}

func decodeSmallTuple(data []byte) (ErlTerm, error) {
	if len(data) < 1 {
		return nil, errors.New("insufficient data for small tuple arity")
	}
	arity := int(data[0])
	terms := make(Tuple, arity)
	remaining := data[1:]
	for i := 0; i < arity; i++ {
		term, bytesRead, err := decodeTerm(remaining)
		if err != nil {
			return nil, err
		}
		terms[i] = term
		remaining = remaining[bytesRead:]
	}
	return Tuple(terms), nil
}

func decodeTerm(data []byte) (ErlTerm, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data to decode")
	}

	typeTag := data[0]
	switch typeTag {
	case ETF_ATOM:
		atom, err := decodeAtom(data[1:])
		if err != nil {
			return nil, 0, err
		}
		return atom, 2 + len(atom.(Atom)), err
	case ETF_SMALL_INTEGER:
		if len(data) < 2 {
			return nil, 0, errors.New("insufficient data for small integer")
		}
		value := Integer(data[1])
		return value, 2, nil
	case ETF_INTEGER:
		if len(data) < 5 {
			return nil, 0, errors.New("insufficient data for integer")
		}
		value := int32(binary.BigEndian.Uint32(data[1:5]))
		return Integer(value), 5, nil
	case ETF_SMALL_TUPLE:
		if len(data) < 2 {
			return nil, 0, errors.New("insufficient data for small tuple arity")
		}
		arity := int(data[1])
		terms := make(Tuple, arity)
		remaining := data[2:]
		totalBytes := 2
		for i := 0; i < arity; i++ {
			term, bytesRead, err := decodeTerm(remaining)
			if err != nil {
				return nil, 0, err
			}
			terms[i] = term
			remaining = remaining[bytesRead:]
			totalBytes += bytesRead
		}
		return Tuple(terms), totalBytes, nil
	case ETF_FLOAT:
		floatVal, err := decodeFloat(data[1:])
		return floatVal, 1 + len(data[1:]), err
	default:
		return nil, 0, fmt.Errorf("unsupported type tag: %d", typeTag)
	}
}

//	let mutable res = double 0
//
// match Double.TryParse (getLatin1 xs 26, &res) with
// | true  -> Float res
// | false -> Error "Float?"
func decodeFloat(data []byte) (ErlTerm, error) {
	if len(data) < 8 {
		return nil, errors.New("insufficient data for float")
	}
	value := binary.BigEndian.Uint64(data[:8])
	return Float(math.Float64frombits(value)), nil
}
