# fibonacci
[![Build Status](https://travis-ci.org/Devtorial/fibonacci.svg?branch=master)](https://travis-ci.org/robarchibald/fibonacci) [![Coverage Status](https://coveralls.io/repos/github/Devtorial/fibonacci/badge.svg?branch=master)](https://coveralls.io/github/robarchibald/fibonacci?branch=master)

A simple API server written in Go which will output a user-defined number of values in the Fibonacci sequence (from 0 - 92 after which int64 overflows). Also included is a simple web UI to show it running in a UI. Go to [https://virtustream.endfirst.com](https://virtustream.endfirst.com) for a working example.

### To use:
1. `git clone github.com/robarchibald/fibonacci`
2. `cd fibonacci/server`
3. `go run fibonacci.go` (note: if you run from a folder other than the server folder, use the -c flag to set the correct client HTML folder)
4. Open a web browser to `http://localhost:1123` to see the web UI or access the api directly at `http://localhost:1123/api/fib/{number from 0-92}`

### Command-line options:
- `l`: Log directory. Default: /var/log
- `c`: Client html directory. Default: ../client
- `a`: Listen Address for server. Default: localhost
- `p`: Listen Port for server. Default: 1123

### Production deployment notes:
- Server can be installed on a system using systemd with the configuration file located at `systemd/virtustream.service`
  - Build server using `go build` from the server folder
  - Move `server` to `/usr/bin` with `mv server /usr/bin/virtustream`
  - Copy systemd init script to `/lib/systemd/system` with `cp ../systemd/virtustream.service /lib/systemd/system`
  - Set ownership to root with `chown root:root /lib/systemd/system/virtustream.service`
  - Set permission to 644 with `chmod 644 /lib/systemd/system/virtustream.service`
  - Enable the service with `systemctl enable virtustream`
  - Reload configs with `systemctl daemon-reload`
  - Start the service with `systemctl start virtustream`
- Alternatively, a Docker image could be created to host the API server executable using the `scratch` image. In order to do that, we'll need to make a slight change to the code. We'll need to turn off logging completely and instead just utilize the built-in logging in NGINX. We also need to recompile into a static binary.
  - Remove logging from fibonacci.go by changing the handler to `handler := handlers.CompressHandler(router)`
  - Compile as a static binary without CGO enabled using `CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .`
  - Copy the Dockerfile to the current directory `cp ../docker/Dockerfile .`
  - Build the Docker image with: `docker build -t virtustreamAPI -f Dockerfile .`
  - Run the image with `docker run -p 127.0.0.1:1123:1123 --name virtustream -t virtustreamAPI`
- Server can be installed behind an NGINX reverse proxy server using the configuration file in `nginx/virtustream.endfirst.com.conf`
  - Copy the configuration file using `cp ../nginx/nginx/virtustream.endfirst.com.conf /etc/nginx/sites-enabled`
  - Create the web folder using `mkdir /usr/share/nginx/virtustreamHtml`
  - Copy web files using `cp ../client/* /usr/share/nginx/virtustreamHtml`
  - Restart NGINX using `service nginx restart`

