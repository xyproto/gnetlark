# <img src="img/gnetlark.png" width=48 alt="gnetlark logo"> gnetlark

[![Build Status](https://travis-ci.org/xyproto/gnetlark.svg?branch=master)](https://travis-ci.org/xyproto/gnetlark)

Simple HTTP server that supports handlers written in Starlark.

The functionality available in the Starlark handlers is currently very limited, but it's fast and secure.

`gnetlark` offers an easy way to try out the Starlark programming language (which is very similar to Python).

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

* [gnet](https://github.com/panjf2000/gnet), for serving HTTP.
* [starlark-go](https://github.com/google/starlark-go), for running Starlark scripts.

### Screenshot

Screenshot of a page served by [`index.star`](index.star), with the server running on port `7711`:

![screenshot](img/screenshot.png)

### Example

A short Starlark script for handling requests and outputting "Hello, World!" ([`hello.star`](hello.star)):

```python
def index(status, msg, method, path, date):
    return "HTTP/1.1 " + status + "\r\nServer: gnetlark\r\nDate: " + date + "\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n" + "Hello, World!"
```

## General info

* Version: 0.0.1
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
* License: MIT
