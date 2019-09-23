package gnetlark

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/xyproto/textoutput"
	"go.starlark.net/starlark"
)

// Request contains fields related to a HTTP request, such as HTTP method
type Request struct {
	Proto      string
	Method     string
	Path       string
	Query      string
	Head       string
	Body       string
	RemoteAddr string
}

// Respond will append a valid HTTP response to the provide bytes.
// The status param should be the code plus text such as "200 OK".
func Respond(to *textoutput.TextOutput, b []byte, sourceFilename, handlerName, status, msg, method, path string) []byte {
	thread := &starlark.Thread{Name: "a thread"}
	globals, err := starlark.ExecFile(thread, sourceFilename, nil, nil)
	if err != nil {
		to.Err("error: " + err.Error())
		return []byte{}
	}
	handlerFunc, ok := globals[handlerName]
	if !ok {
		to.Err("error: could not find function " + to.LightBlue(handlerName) + " in " + to.LightGreen(sourceFilename))
		return []byte{}
	}
	datestring := time.Now().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	v, err := starlark.Call(thread, handlerFunc, starlark.Tuple{starlark.String(status), starlark.String(msg), starlark.String(method), starlark.String(path), starlark.String(datestring)}, nil)
	if err != nil {
		to.Err("error: " + err.Error())
		return []byte{}
	}
	starString, ok := v.(starlark.String)
	if !ok {
		to.Err("error: " + to.LightBlue(handlerName) + " in " + to.LightGreen(sourceFilename) + " returned something that was not a string")
		return []byte{}
	}
	//log.Println("returning " + starString.GoString())
	return []byte(starString.GoString())
}

// ParseReq is a very simple http request parser. This operation
// waits for the entire payload to be buffered before returning a
// valid request.
func ParseReq(data []byte, req *Request) (leftover []byte, err error) {
	sdata := string(data)
	var i, s int
	var top string
	var clen int
	var q = -1
	// method, path, proto line
	for ; i < len(sdata); i++ {
		if sdata[i] == ' ' {
			req.Method = sdata[s:i]
			for i, s = i+1, i+1; i < len(sdata); i++ {
				if sdata[i] == '?' && q == -1 {
					q = i - s
				} else if sdata[i] == ' ' {
					if q != -1 {
						req.Path = sdata[s:q]
						req.Query = req.Path[q+1 : i]
					} else {
						req.Path = sdata[s:i]
					}
					for i, s = i+1, i+1; i < len(sdata); i++ {
						if sdata[i] == '\n' && sdata[i-1] == '\r' {
							req.Proto = sdata[s:i]
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
	if req.Proto == "" {
		return data, errors.New("malformed request")
	}
	top = sdata[:s]
	for ; i < len(sdata); i++ {
		if i > 1 && sdata[i] == '\n' && sdata[i-1] == '\r' {
			line := sdata[s : i-1]
			s = i + 1
			if line == "" {
				req.Head = sdata[len(top)+2 : i+1]
				i++
				if clen > 0 {
					if len(sdata[i:]) < clen {
						break
					}
					req.Body = sdata[i : i+clen]
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
