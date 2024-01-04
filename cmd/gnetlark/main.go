package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/panjf2000/gnet/v2"
	"github.com/xyproto/gnetlark"
	"github.com/xyproto/textoutput"
)

var res string

type httpServer struct {
	gnet.BuiltinEventEngine

	port           int
	multicore      bool
	to             *textoutput.TextOutput
	sourceFilename string
}

func (hs *httpServer) OnBoot(_ gnet.Engine) gnet.Action {
	hs.to.OutputTags(fmt.Sprintf("<lightgreen>HTTP server is listening on port %d<off>\n", hs.port))
	return gnet.None
}

func (hs *httpServer) OnOpen(_ gnet.Conn) ([]byte, gnet.Action) {
	return nil, gnet.None
}

func (hs *httpServer) OnTraffic(c gnet.Conn) gnet.Action {
	buffer := make([]byte, 1024) // Adjust the buffer size as needed
	n, err := c.Read(buffer)
	if err != nil {
		hs.to.OutputTags("<red>Error reading data:</red> " + err.Error())
		return gnet.Close
	}

	data := buffer[:n]

	var req gnetlark.Request
	leftover, err := gnetlark.ParseReq(data, &req)
	if err != nil {
		hs.to.OutputTags("<red>Server error:</red> " + err.Error())
		responseBytes := gnetlark.Respond(hs.to, nil, hs.sourceFilename, "error", "500 Error", err.Error()+"\n", req.Method, req.Path)
		c.Write(responseBytes)
		return gnet.Close
	} else if len(leftover) == len(data) {
		hs.to.OutputTags("<yellow>Request not ready<off>")
		return gnet.None
	}

	req.RemoteAddr = c.RemoteAddr().String()
	responseBytes := gnetlark.Respond(hs.to, nil, hs.sourceFilename, "index", "200 OK", res, req.Method, req.Path)
	c.Write(responseBytes)
	return gnet.None
}

func main() {
	var (
		port           int
		multicore      bool
		sourceFilename string
	)

	flag.IntVar(&port, "port", 80, "server port")
	flag.BoolVar(&multicore, "multicore", true, "multicore")
	flag.StringVar(&sourceFilename, "main", "index.star", "main script")
	flag.Parse()

	to := textoutput.New()

	hs := &httpServer{
		port:           port,
		multicore:      multicore,
		to:             to,
		sourceFilename: sourceFilename,
	}

	addr := fmt.Sprintf("tcp://:%d", port)
	to.OutputTags("<lightyellow>Server starting...<off>")
	log.Println("Server exits:", gnet.Run(hs, addr, gnet.WithMulticore(multicore)))
}
