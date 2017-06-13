package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
)

var errInvalidFormat = "Invalid url format /api/fib/{number} where {number} is an integer\nExample: /api/fib/5 to retrieve the first 5 digits in the fibonacci sequence"
var errInvalidNumber = "Input must be between 0 and 92"

type server interface {
	ListenAndServe() error
}

var creator = func(addr string, handler http.Handler) server {
	return &http.Server{Addr: addr, Handler: handler}
}

func main() {
	logDir, clientHTMLDir, listenAt := getFlags()
	logFile := openLog(logDir)
	defer logFile.Close()
	server := getServer(listenAt, clientHTMLDir, logFile)
	server.ListenAndServe()
}

func getFlags() (string, string, string) {
	logDir := flag.String("l", "/var/log", "Log directory. Default: /var/log")
	clientHTMLDir := flag.String("c", "../client", "Client html directory. Default: ../client")
	listenAddress := flag.String("a", "", "Listen Address for server. Default: localhost")
	listenPort := flag.Int("p", 1123, "Listen Port for server. Default: 1123")
	flag.Parse()
	listenAt := fmt.Sprintf("%s:%d", *listenAddress, *listenPort)
	return *logDir, *clientHTMLDir, listenAt
}

func openLog(logDir string) *os.File {
	logFile, err := os.OpenFile(path.Join(logDir, "fibonacci.log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0640)
	if err != nil {
		log.Fatal("unable to open log file", err)
	}
	return logFile
}

func getServer(listenAt, clientHTMLDir string, logFile *os.File) server {
	router := httprouter.New()
	router.GET("/api/fib/:number", getFib)                                                // fibonacci API server
	router.NotFound = http.FileServer(http.Dir(clientHTMLDir))                            // html client
	handler := handlers.CombinedLoggingHandler(logFile, handlers.CompressHandler(router)) // handle compression, routing and logging
	return creator(listenAt, handler)
}

// Taken from golang.org fibonacci closure example
func fib() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func getFibSequence(num int) ([]int, error) {
	if num < 0 || num > 92 {
		return nil, errors.New(errInvalidNumber) // negative numbers are invalid and > 92 will overflow the integer
	}
	f := fib()
	r := make([]int, num)
	for i := 0; i < num; i++ {
		r[i] = f()
	}
	return r, nil
}

func getFib(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	num, err := strconv.Atoi(p.ByName("number"))
	if err != nil {
		http.Error(w, errInvalidFormat, http.StatusBadRequest)
		return
	}

	sequence, err := getFibSequence(num)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // accept between 0-92 so it is still a client error
		return
	}

	json, err := json.Marshal(sequence)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // unable to encode json so server error
		return
	}

	fmt.Fprintf(w, string(json))
}
