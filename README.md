# <img src="img/gnetlark.png" width=48 alt="gnetlark logo"> gnetlark

[![Build Status](https://travis-ci.org/xyproto/gnetlark.svg?branch=master)](https://travis-ci.org/xyproto/gnetlark)

Simple HTTP server that supports handlers written in Starlark.

The functionality available in the Starlark handlers is currently very limited, but it's fast, secure and a possibility to try out the Starlark language.

### Installation

One way of building and installing

    git clone https://github.com/xyproto/gnetlark
    cd gnetlark
    go build
    sudo install -Dm755 gnetlark /usr/bin/gnetlark

### Configuration

One way to allow access to port 80 on Linux:

    sudo setcap cap_net_bind_service=+ep /usr/bin/gnetlark

It's also possible to specify a port with `--port` or run it as root (not recommended).

### Depends on

* [gnet](https://github.com/panjf2000/gnet), for serving HTTP
* [starlark-go](https://github.com/google/starlark-go), for running Starlark scripts

### Screenshot

Screenshot of a page served by `index.star`, with the server running on port 7711:

![screenshot](img/screenshot.png)

## General info

* Version: 0.0.1
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
* License: MIT
