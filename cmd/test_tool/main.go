package main

import (
	"etf"
	"fmt"
)

func main() {
	// Example term: {atom, 42}
	term := etf.Tuple{
		etf.Atom("atom"),
		etf.Integer(42),
		etf.Float(3.14),
	}

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
