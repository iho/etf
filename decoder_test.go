package etf_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/iho/etf"
)

func invokeErl(term string) ([]byte, error) {
	erlExpr := fmt.Sprintf(`
        erlang:display(term_to_binary(
            "%s"
        )), halt().`, term)

	cmd := exec.Command("erl", "-noinput", "-eval", erlExpr, "-s", "init", "stop")
	// Erlang Term: <<131,107,0,4,97,116,111,109>>
	res, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	resStr := string(res)

	fmt.Println("Erlang Term:", resStr)
	resStr = strings.TrimPrefix(resStr, "<<")
	fmt.Println("Erlang Term:", resStr)
	resStr = strings.TrimSuffix(resStr, ">>")
	fmt.Println("Erlang Term:", resStr)

	resArray := strings.Split(resStr, ",")
	resArr := make([]byte, len(resArray))
	for i, v := range resArray {
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		resArr[i] = byte(val)
	}
	return resArr, nil
}

func ExampleEncodeErlTerm() {
	a, b := etf.Atom("a"), etf.Atom("b")
	term := etf.Tuple{
		etf.Integer(999),
		etf.Atom("atom"),
		etf.Nil{},
		etf.List{a, b, a, b},
		etf.Integer(10),
		etf.Binary{0x01, 0x02, 0x03},
		etf.Float(3.14),
		etf.Integer(42),
		etf.Map{
			"key1": etf.Atom("value1"),
			"key2": etf.Atom("value2"),
		},
	}
	_, err := etf.EncodeErlTerm(term, true)
	if err != nil {
		fmt.Println("Encoding error:", err)
		return
	}
}

func TestErAtomlTerm(t *testing.T) {
	term := etf.Atom("atom")
	encoded, err := etf.EncodeErlTerm(term, true)
	if err != nil {
		t.Error("Encoding error:", err)
	}

	erlangTerm, err := invokeErl("atom")
	if err != nil {
		t.Error("Encoding error:", err)
	}

	fmt.Println("Encoded ETF Binary:", encoded)
	fmt.Println("Erlang Term:", string(erlangTerm))
	if !bytes.Equal(encoded, erlangTerm) {
		t.Errorf("Expected %v, got %v", erlangTerm, encoded)
	}
}
