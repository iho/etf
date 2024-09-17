package main

import (
	"fmt"

	"github.com/iho/etf"
)

func main() {
	// a, b := etf.Atom("a"), etf.Atom("b")
	term := etf.Tuple{
		// etf.Integer(999),
		// etf.Atom("atom"),
		// etf.Nil{},
		// etf.List{a, b, a, b},
		// etf.Integer(10),
		// etf.Binary{0x01, 0x02, 0x03},
		// etf.Float(3.14),

		etf.Map{
			// etf.MapElem{Key: etf.Atom("key1"), Value: etf.Atom("value1")},
			// etf.MapElem{Key: etf.Atom("key2"), Value: etf.Atom("value2")},
			etf.MapElem{Key: etf.Tuple{etf.Atom("key2"), etf.Integer(10)}, Value: etf.Atom("value2")},
		},
		// etf.Integer(42),
	}

	fmt.Printf("Term: %+v\n", term)

	// Encode the term
	encoded, err := etf.EncodeErlTerm(term, true)
	if err != nil {
		fmt.Println("Encoding error:", err)
		return
	}

	fmt.Printf("Encoded ETF Binary: %v\n", encoded)

	// Decode the term
	decoded, err := etf.DecodeErlTerm(encoded)
	if err != nil {
		fmt.Println("Decoding error:", err)
		return
	}

	fmt.Printf("Decoded Term: %+v\n", decoded)
}
