# gnetlark

A fast HTTP server that supports simple handlers written in Starlark.

## Uses these go packages

* [gnet](https://github.com/panjf2000/gnet)
* [starlark-go](https://github.com/google/starlark-go)

## Allow access to port 80 on Linux

    sudo setcap cap_net_bind_service=+ep /usr/bin/gnetlark

It's also possible to specify a port with `--port` or run it as root (not recommended).

## General info

* License: MIT
* Version: 0.0.1
* Author: Alexander F. Rødseth &lt;xyproto@archlinux.org&gt;
