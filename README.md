# fibonacci
[![Build Status](https://travis-ci.org/robarchibald/fibonacci.svg?branch=master)](https://travis-ci.org/robarchibald/fibonacci) [![Coverage Status](https://coveralls.io/repos/github/robarchibald/fibonacci/badge.svg?branch=master)](https://coveralls.io/github/robarchibald/fibonacci?branch=master)

A simple API server written in Go which will output a user-defined number of values in the Fibonacci sequence (from 0 - 92 after which int64 overflows). Also included is a simple web UI to show it running in a UI.

To use:
1. `git clone github.com/robarchibald/fibonacci`
2. `cd fibonacci/server`
3. `go run fibonacci.go` (note: if you run from a folder other than the server folder, use the -c flag to set the correct client HTML folder)
4. Open a web browser to `http://localhost:1123` to see the web UI or access the api directly at `http://localhost:1123/api/fib/{number from 0-92}`

Command-line options:
- `l`: Log directory. Default: /var/log
- `c`: Client html directory. Default: ../client
- `a`: Listen Address for server. Default: localhost
- `p`: Listen Port for server. Default: 1123


