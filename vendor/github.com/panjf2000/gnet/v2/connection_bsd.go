// Copyright (c) 2021 The Gnet Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build freebsd || dragonfly || netbsd || openbsd || darwin
// +build freebsd dragonfly netbsd openbsd darwin

package gnet

import (
	"golang.org/x/sys/unix"

	"github.com/panjf2000/gnet/v2/internal/netpoll"
)

func (c *conn) handleEvents(_ int, filter int16) (err error) {
	switch filter {
	case netpoll.EVFilterSock:
		err = c.loop.close(c, unix.ECONNRESET)
	case netpoll.EVFilterWrite:
		if !c.outboundBuffer.IsEmpty() {
			err = c.loop.write(c)
		}
	case netpoll.EVFilterRead:
		err = c.loop.read(c)
	}
	return
}
