package etf

type ErlTerm any

type Atom string

type Integer int16

type Float float64

type Tuple []ErlTerm

type MapElem struct {
	Key   ErlTerm
	Value ErlTerm
}

type Map []MapElem

type List []ErlTerm

type Binary []byte

type Nil struct{}
