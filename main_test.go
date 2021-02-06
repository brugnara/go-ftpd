package main

import (
	"errors"
	"log"
	"net"
	"testing"
	"time"
)

// io.ReadWriteCanceler interface

type rwc struct {
	in       []string
	out      []string
	isClosed bool
}

func (r *rwc) Close() error {
	if r.isClosed {
		return errors.New("already closed")
	}
	r.isClosed = true
	return nil
}

func (r *rwc) Write(b []byte) (int, error) {
	r.out = append(r.out, string(b))
	log.Printf("Received: %d bytes\n", len(b))
	return len(b), nil
}

func (r *rwc) Read(b []byte) (int, error) {
	if len(r.in) == 0 {
		return 0, nil
	}
	// io.Copy(bytes.NewBuffer(b), strings.NewReader(r.in[0]))
	// b = []byte(r.in[0])
	copy(b, r.in[0])
	ln := len(r.in[0])
	log.Printf("Sent: %d bytes\n", ln)
	r.in = r.in[1:]
	return ln, nil
}

func (r *rwc) LocalAddr() net.Addr                { return nil }
func (r *rwc) RemoteAddr() net.Addr               { return nil }
func (r *rwc) SetDeadline(t time.Time) error      { return nil }
func (r *rwc) SetReadDeadline(t time.Time) error  { return nil }
func (r *rwc) SetWriteDeadline(t time.Time) error { return nil }

// net.Listener interface
type lstn struct {
	accept   int
	accepted int
	closes   int
}

func (l *lstn) Accept() (net.Conn, error) {
	if l.accepted >= l.accept {
		return nil, errors.New("No more accept")
	}
	l.accepted++
	return &rwc{}, nil
}

func (l *lstn) Close() error {
	l.closes++
	return nil
}

func (l *lstn) Addr() net.Addr { return nil }

// test stuff

func TestHandlerCloser(t *testing.T) {
	conn := &rwc{}
	handler(conn)
	if !conn.isClosed {
		t.Error("Conn was not closed!")
	}
}

func TestHandler(t *testing.T) {
	for _, test := range []struct {
		in  []string
		out int
	}{
		{[]string{"ls"}, 29},
		{[]string{"cd"}, 5},
		{[]string{"cd", "cd"}, 5},
		{[]string{"cd tt"}, 4},
		{[]string{"cd tt", "ls"}, 5},
		{[]string{"cddasd"}, 5},
	} {
		conn := &rwc{in: test.in}
		handler(conn)
		if o := conn.out; len(o) != test.out {
			t.Error("Expected:", len(o), "to be:", test.out, "with:", test.in)
		}
	}
}

func TestValidPath(t *testing.T) {
	for _, test := range []struct {
		in  string
		out bool
	}{
		{"./public", true},
		{"public", true},
		{"foo", false},
		{"", false},
	} {
		if o := validPath(test.in); o != test.out {
			t.Error("Expected:", o, "to be:", test.out, "with:", test.in)
		}
	}
}

func TestLoop(t *testing.T) {
	stopped := true
	listener := &lstn{accept: 1}
	loop(listener, &stopped)
	//
	if listener.accepted != 1 {
		t.Error("Invalid accepted count")
	}
	if listener.closes != 1 {
		t.Error("Wrong closes count")
	}
}
