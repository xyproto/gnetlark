package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/panjf2000/gnet"
	"github.com/xyproto/textoutput"
	"go.starlark.net/starlark"
)

var res string

type request struct {
	proto      string
	method     string
	path       string
	query      string
	head       string
	body       string
	remoteAddr string
}

func main() {
	var (
		port   int
		colors bool
		quiet  bool
	)

	flag.IntVar(&port, "port", 80, "server port")
	flag.BoolVar(&colors, "colors", true, "enable colors")
	flag.BoolVar(&quiet, "quiet", false, "no output")

	flag.Parse()

	to := textoutput.NewTextOutput(colors, !quiet)

	var events gnet.Events
	events.Multicore = true
	events.OnInitComplete = func(srv gnet.Server) (action gnet.Action) {
		log.Printf("HTTP server started on port %d", port)
		return
	}

	events.React = func(c gnet.Conn) (out []byte, action gnet.Action) {
		top, tail := c.ReadPair()
		data := append(top, tail...)
		var req request
		leftover, err := parsereq(data, &req)
		if err != nil {
			log.Println("Server error: " + err.Error())
			out = appendresp(to, out, "index.star", "error", "500 Error", err.Error()+"\n")
			action = gnet.Close
			return
		} else if len(leftover) == len(data) {
			log.Println("Request not ready")
			return
		}

		// handle the request
		req.remoteAddr = c.RemoteAddr().String()
		out = appendresp(to, out, "index.star", "index", "200 OK", res)
		c.ResetBuffer()
		return
	}
	// We at least want the single http address.
	addrs := []string{fmt.Sprintf("tcp"+"://:%d", port)}
	// Start serving!
	log.Fatal(gnet.Serve(events, addrs...))
}

// appendresp will append a valid http response to the provide bytes.
// The status param should be the code plus text such as "200 OK".
func appendresp(to *textoutput.TextOutput, b []byte, sourceFilename, handlerName, status, msg string) []byte {
	thread := &starlark.Thread{Name: "a thread"}
	globals, err := starlark.ExecFile(thread, sourceFilename, nil, nil)
	if err != nil {
		to.Err("error: " + err.Error())
	}
	handlerFunc, ok := globals[handlerName]
	if !ok {
		to.Err("error: could not find function " + to.LightBlue(handlerName) + " in " + to.LightGreen(sourceFilename))
	}
	datestring := time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	v, err := starlark.Call(thread, handlerFunc, starlark.Tuple{starlark.String(status), starlark.String(msg), starlark.String(datestring)}, nil)
	if err != nil {
		to.Err("error: " + err.Error())
	}
	starString, ok := v.(starlark.String)
	if !ok {
		to.Err("error: " + to.LightBlue(handlerName) + " in " + to.LightGreen(sourceFilename) + " returned something that was not a string")
	}
	//log.Println("returning " + starString.GoString())
	return []byte(starString.GoString())
}

// parsereq is a very simple http request parser. This operation
// waits for the entire payload to be buffered before returning a
// valid request.
func parsereq(data []byte, req *request) (leftover []byte, err error) {
	sdata := string(data)
	var i, s int
	var top string
	var clen int
	var q = -1
	// method, path, proto line
	for ; i < len(sdata); i++ {
		if sdata[i] == ' ' {
			req.method = sdata[s:i]
			for i, s = i+1, i+1; i < len(sdata); i++ {
				if sdata[i] == '?' && q == -1 {
					q = i - s
				} else if sdata[i] == ' ' {
					if q != -1 {
						req.path = sdata[s:q]
						req.query = req.path[q+1 : i]
					} else {
						req.path = sdata[s:i]
					}
					for i, s = i+1, i+1; i < len(sdata); i++ {
						if sdata[i] == '\n' && sdata[i-1] == '\r' {
							req.proto = sdata[s:i]
							i, s = i+1, i+1
							break
						}
					}
					break
				}
			}
			break
		}
	}
	if req.proto == "" {
		return data, fmt.Errorf("malformed request")
	}
	top = sdata[:s]
	for ; i < len(sdata); i++ {
		if i > 1 && sdata[i] == '\n' && sdata[i-1] == '\r' {
			line := sdata[s : i-1]
			s = i + 1
			if line == "" {
				req.head = sdata[len(top)+2 : i+1]
				i++
				if clen > 0 {
					if len(sdata[i:]) < clen {
						break
					}
					req.body = sdata[i : i+clen]
					i += clen
				}
				return data[i:], nil
			}
			if strings.HasPrefix(line, "Content-Length:") {
				n, err := strconv.ParseInt(strings.TrimSpace(line[len("Content-Length:"):]), 10, 64)
				if err == nil {
					clen = int(n)
				}
			}
		}
	}
	// not enough data
	return data, nil
}
