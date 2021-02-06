package main

import (
	"errors"
	"log"
	"testing"
)

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
