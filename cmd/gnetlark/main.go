package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/panjf2000/gnet"
	"github.com/xyproto/gnetlark"
	"github.com/xyproto/textoutput"
)

var res string

func main() {
	var (
		port           int
		colors         bool
		quiet          bool
		sourceFilename string
	)

	flag.IntVar(&port, "port", 80, "server port")
	flag.BoolVar(&colors, "colors", true, "enable colors")
	flag.BoolVar(&quiet, "quiet", false, "no output")
	flag.StringVar(&sourceFilename, "main", "index.star", "main script")

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
		var req gnetlark.Request
		leftover, err := gnetlark.ParseReq(data, &req)
		if err != nil {
			log.Println("Server error: " + err.Error())
			out = gnetlark.Respond(to, out, sourceFilename, "error", "500 Error", err.Error()+"\n", req.Method, req.Path)
			action = gnet.Close
			return
		} else if len(leftover) == len(data) {
			log.Println("Request not ready")
			return
		}

		// handle the request
		req.RemoteAddr = c.RemoteAddr().String()
		out = gnetlark.Respond(to, out, sourceFilename, "index", "200 OK", res, req.Method, req.Path)
		c.ResetBuffer()
		return
	}
	// We at least want the single HTTP address.
	addrs := []string{fmt.Sprintf("tcp"+"://:%d", port)}

	// Serve forever, but quit with an error if it stops
	to.ErrExit(gnet.Serve(events, addrs...).Error())
}
