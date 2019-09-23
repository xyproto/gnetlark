// Copyright 2019 Andy Pan. All rights reserved.
// Copyright 2018 Joshua J Baker. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// +build darwin netbsd freebsd openbsd dragonfly linux

package gnet

import (
	"net"
	"sync"
	"time"

	"github.com/panjf2000/gnet/netpoll"
	"github.com/panjf2000/gnet/ringbuffer"
	"golang.org/x/sys/unix"
)

type loop struct {
	idx         int             // loop index in the server loops list
	poller      *netpoll.Poller // epoll or kqueue
	packet      []byte          // read packet buffer
	connections sync.Map        // loop connections fd -> conn
}

func (l *loop) loopCloseConn(svr *server, conn *conn, err error) error {
	l.poller.Delete(conn.fd)
	//delete(l.connections, conn.fd)
	l.connections.Delete(conn.fd)
	_ = unix.Close(conn.fd)

	if svr.events.OnClosed != nil {
		switch svr.events.OnClosed(conn, err) {
		case None:
		case Shutdown:
			return ErrClosing
		}
	}
	return nil
}

func (l *loop) loopNote(svr *server, note interface{}) error {
	var err error
	switch v := note.(type) {
	case time.Duration:
		delay, action := svr.events.Tick()
		switch action {
		case None:
		case Shutdown:
			err = ErrClosing
		}
		svr.tch <- delay
	case error: // shutdown
		err = v
	case *conn:
		// Wake called for connection
		if val, ok := l.connections.Load(v.fd); !ok || val != v {
			return nil // ignore stale wakes
		}
		return l.loopWake(svr, v)
	//case *mail:
	//	l.fdconns[v.fd] = v.conn
	//	_ = l.loopOpened(svr, v.conn)
	//	l.poller.AddRead(v.fd)
	case func():
		v()
	}
	return err
}

func (l *loop) loopRun(svr *server) {
	defer func() {
		svr.signalShutdown()
		svr.wg.Done()
	}()

	if l.idx == 0 && svr.events.Tick != nil {
		go l.loopTicker(svr)
	}

	_ = l.poller.Polling(func(fd int, note interface{}) error {
		if fd == 0 {
			return l.loopNote(svr, note)
		}
		if v, ok := l.connections.Load(fd); ok {
			c := v.(*conn)
			switch {
			case !c.opened:
				return l.loopOpened(svr, c)
			case c.outBuf.Length() > 0:
				return l.loopWrite(svr, c)
			default:
				return l.loopRead(svr, c)
			}
		} else {
			return l.loopAccept(svr, fd)
		}
	})
}

func (l *loop) loopTicker(svr *server) {
	for {
		if err := l.poller.Trigger(time.Duration(0)); err != nil {
			break
		}
		time.Sleep(<-svr.tch)
	}
}

func (l *loop) loopAccept(svr *server, fd int) error {
	for i, ln := range svr.lns {
		if ln.fd == fd {
			if ln.pconn != nil {
				return l.loopUDPRead(svr, i, fd)
			}
			nfd, sa, err := unix.Accept(fd)
			if err != nil {
				if err == unix.EAGAIN {
					return nil
				}
				return err
			}
			if err := unix.SetNonblock(nfd, true); err != nil {
				return err
			}
			conn := &conn{fd: nfd,
				sa:     sa,
				lnidx:  i,
				inBuf:  ringbuffer.New(cacheRingBufferSize),
				outBuf: ringbuffer.New(cacheRingBufferSize),
				loop:   l,
			}
			//l.connections[conn.fd] = conn
			l.connections.Store(conn.fd, conn)
			l.poller.AddReadWrite(conn.fd)
			return nil
		}
	}
	return nil
}

func (l *loop) loopUDPRead(svr *server, lnidx, fd int) error {
	n, sa, err := unix.Recvfrom(fd, l.packet, 0)
	if err != nil || n == 0 {
		return nil
	}
	if svr.events.React != nil {
		var sa6 unix.SockaddrInet6
		switch sa := sa.(type) {
		case *unix.SockaddrInet4:
			sa6.ZoneId = 0
			sa6.Port = sa.Port
			for i := 0; i < 12; i++ {
				sa6.Addr[i] = 0
			}
			sa6.Addr[12] = sa.Addr[0]
			sa6.Addr[13] = sa.Addr[1]
			sa6.Addr[14] = sa.Addr[2]
			sa6.Addr[15] = sa.Addr[3]
		case *unix.SockaddrInet6:
			sa6 = *sa
		}
		conn := &conn{
			addrIndex:  lnidx,
			localAddr:  svr.lns[lnidx].lnaddr,
			remoteAddr: netpoll.SockaddrToUDPAddr(&sa6),
			inBuf:      ringbuffer.New(cacheRingBufferSize),
		}
		_, _ = conn.inBuf.Write(l.packet[:n])
		out, action := svr.events.React(conn)
		if len(out) > 0 {
			if svr.events.PreWrite != nil {
				svr.events.PreWrite()
			}
			sniffError(unix.Sendto(fd, out, 0, sa))
		}
		switch action {
		case Shutdown:
			return ErrClosing
		}
	}
	return nil
}

func (l *loop) loopOpened(svr *server, conn *conn) error {
	conn.opened = true
	conn.addrIndex = conn.lnidx
	conn.localAddr = svr.lns[conn.lnidx].lnaddr
	conn.remoteAddr = netpoll.SockaddrToTCPOrUnixAddr(conn.sa)
	if svr.events.OnOpened != nil {
		out, opts, action := svr.events.OnOpened(conn)
		conn.action = action
		if opts.TCPKeepAlive > 0 {
			if _, ok := svr.lns[conn.lnidx].ln.(*net.TCPListener); ok {
				sniffError(netpoll.SetKeepAlive(conn.fd, int(opts.TCPKeepAlive/time.Second)))
			}
		}

		if len(out) > 0 {
			conn.open(out)
		}
	}
	if conn.outBuf.Length() != 0 {
		l.poller.AddWrite(conn.fd)
	}
	return l.handleAction(svr, conn)
}

func (l *loop) loopWrite(svr *server, conn *conn) error {
	if svr.events.PreWrite != nil {
		svr.events.PreWrite()
	}

	out := conn.outBuf.Bytes()
	n, err := unix.Write(conn.fd, out)
	if err != nil {
		if err == unix.EAGAIN {
			return nil
		}
		return l.loopCloseConn(svr, conn, err)
	}
	conn.outBuf.Advance(n)
	ringbuffer.Recycle(out)

	if conn.outBuf.Length() == 0 {
		l.poller.ModRead(conn.fd)
	}
	return nil
}

func (l *loop) handleAction(svr *server, conn *conn) error {
	switch conn.action {
	case None:
		return nil
	case Close:
		return l.loopCloseConn(svr, conn, nil)
	case Shutdown:
		return ErrClosing
	default:
		return nil
	}
}

func (l *loop) loopWake(svr *server, conn *conn) error {
	out, action := svr.events.React(conn)
	conn.action = action
	if len(out) > 0 {
		conn.write(out)
	}
	return l.handleAction(svr, conn)
}

func (l *loop) loopRead(svr *server, conn *conn) error {
	n, err := unix.Read(conn.fd, l.packet)
	if n == 0 || err != nil {
		if err == unix.EAGAIN {
			return nil
		}
		return l.loopCloseConn(svr, conn, err)
	}

	_, _ = conn.inBuf.Write(l.packet[:n])
	out, action := svr.events.React(conn)
	conn.action = action
	if len(out) > 0 {
		conn.write(out)
	}
	return l.handleAction(svr, conn)
}
