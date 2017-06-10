package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

const errInvalidFormat = "Invalid url format /api/fib/{number} where {number} is an integer\nExample: /api/fib/5 to retrieve the first 5 digits in the fibonacci sequence"
const errInvalidNumber = "Input must be between 0 and 92"

// Taken from golang.org fibonacci closure example
func fib() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func getFibSequence(num int) ([]int, error) {
	f := fib()
	if num < 0 || num > 92 {
		return nil, errors.New(errInvalidNumber) // negative numbers are invalid and > 92 will overflow the integer
	}
	r := make([]int, num)
	for i := 0; i < num; i++ {
		r[i] = f()
	}
	return r, nil
}

func getFib(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	num, err := strconv.Atoi(p.ByName("number"))
	if err != nil {
		http.Error(w, errInvalidNumber, http.StatusInternalServerError)
		return
	}
	seq, err := getFibSequence(num)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(seq)
	if err != nil {
		http.Error(w, errInvalidNumber, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(json))
}

func main() {
	router := httprouter.New()
	router.GET("/api/fib/:number", getFib)
	log.Fatal(http.ListenAndServe(":1123", router))
}
