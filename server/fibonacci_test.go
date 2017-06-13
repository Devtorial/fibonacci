package main

import (
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/julienschmidt/httprouter"
)

type fakeServer struct {
	listenCalled bool
}

func (s *fakeServer) ListenAndServe() error {
	s.listenCalled = true
	return nil
}

var expected = []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34, 55,
	89, 144, 233, 377, 610, 987, 1597, 2584, 4181, 6765,
	10946, 17711, 28657, 46368, 75025, 121393, 196418, 317811, 514229, 832040,
	1346269, 2178309, 3524578, 5702887, 9227465, 14930352, 24157817, 39088169, 63245986, 102334155,
	165580141, 267914296, 433494437, 701408733, 1134903170, 1836311903, 2971215073, 4807526976, 7778742049, 12586269025,
	20365011074, 32951280099, 53316291173, 86267571272, 139583862445, 225851433717, 365435296162, 591286729879, 956722026041, 1548008755920,
	2504730781961, 4052739537881, 6557470319842, 10610209857723, 17167680177565, 27777890035288, 44945570212853, 72723460248141, 117669030460994, 190392490709135,
	308061521170129, 498454011879264, 806515533049393, 1304969544928657, 2111485077978050, 3416454622906707, 5527939700884757, 8944394323791464, 14472334024676221, 23416728348467685,
	37889062373143906, 61305790721611591, 99194853094755497, 160500643816367088, 259695496911122585, 420196140727489673, 679891637638612258, 1100087778366101931, 1779979416004714189, 2880067194370816120,
	4660046610375530309, 7540113804746346429, // anything after this number will overflow the int
}

func TestFib(t *testing.T) {
	f := fib()
	// values taken from http://www.maths.surrey.ac.uk/hosted-sites/R.Knott/Fibonacci/fibtable.html#series
	// fib function starts sequence at 1 so skip the first value
	for i := 1; i < len(expected); i++ {
		actual := f()
		if expected[i] != actual {
			t.Errorf("mismatch in index %d. Expected vs. actual: %d|%d\n", i, expected[i], actual)
		}
	}
}

func TestGetFibSequence(t *testing.T) {
	if _, err := getFibSequence(-1); err == nil {
		t.Error("expected out of range error")
	}
	if _, err := getFibSequence(99); err == nil {
		t.Error("expected out of range error")
	}

	// Max allowed value for getFibSequence
	actual, _ := getFibSequence(93)
	if len(actual) != 93 {
		t.Error("expected a sequence of 93, got", len(actual))
	} else {
		for i := 0; i < 93; i++ {
			if expected[i] != actual[i] {
				t.Errorf("getFibSequence(93) expected and actual don't match. Expected vs. actual: %d|%d", expected[i], actual[i])
			}
		}
	}

	// Min allowed value for getFibSequence
	actual, _ = getFibSequence(0)
	if len(actual) != 0 {
		t.Error("expected a sequence of 0, got", len(actual))
	}

	// Random allowed value for getFibSequence
	actual, _ = getFibSequence(5)
	if len(actual) != 5 {
		t.Error("expected a sequence of 5, got", len(actual))
	} else {
		for i := 0; i < 5; i++ {
			if expected[i] != actual[i] {
				t.Errorf("getFibSequence(5) expected and actual don't match. Expected vs. actual: %d|%d", expected[i], actual[i])
			}
		}
	}
}

func TestGetFib(t *testing.T) {
	// Error due to no number being passed in
	w := httptest.NewRecorder()
	getFib(w, nil, []httprouter.Param{})
	if body := w.Body.String(); body != errInvalidFormat+"\n" {
		t.Error("expected invalid format", body)
	}

	// Error due to number out of range
	w = httptest.NewRecorder()
	getFib(w, nil, []httprouter.Param{httprouter.Param{Key: "number", Value: "-3"}})
	if body := w.Body.String(); body != errInvalidNumber+"\n" {
		t.Error("expected invalid number", body)
	}

	// success
	w = httptest.NewRecorder()
	getFib(w, nil, []httprouter.Param{httprouter.Param{Key: "number", Value: "5"}})
	if body := w.Body.String(); body != "[0,1,1,2,3]" {
		t.Error("expected valid result", body)
	}
}

func TestGetFlags(t *testing.T) {
	logDir, clientHTMLDir, listenAt := getFlags()
	if logDir != "/var/log" || clientHTMLDir != "../client" || listenAt != ":1123" {
		t.Error("expected default flag values")
	}
}

func TestGetServer(t *testing.T) {
	expectedAddr := "hello"
	// use fake server, not real
	s := getServer(expectedAddr, "htmldir", nil).(*http.Server)
	if s.Addr != expectedAddr {
		t.Error("expected correct address")
	}
}

func TestMain(t *testing.T) {
	s := &fakeServer{}
	creator = func(addr string, handler http.Handler) server {
		return s
	}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	os.Args = []string{os.Args[0], "-l", "."} // set logDir to current folder

	main()
	if !s.listenCalled {
		t.Error("expected ListenAndServe to be called")
	}
}
