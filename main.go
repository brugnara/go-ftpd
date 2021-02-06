package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
)

const host = "localhost:8021"

var done chan bool
var rootPath string

func init() {
	done = make(chan bool)
	//
	flag.StringVar(&rootPath, "path", "public", "Base path to lock FTP to")
	flag.Parse()
}

func main() {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		panic(err)
	}
	//
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("Exiting now...")
		done <- true
	}()
	//
	log.Printf("Starting FTP with root folder:\n\t> %s\n", rootPath)

	if !validPath(rootPath) {
		panic("Invalid path")
	}

	stop := false

	go loop(listener, &stop)
	stop = <-done
}

func validPath(p string) (valid bool) {
	f, err := os.Open(p)
	if err == nil {
		defer f.Close()
		valid = true
	}
	return
}

func handler(c io.ReadWriteCloser) {
	defer c.Close()
	ftp := newFtp(rootPath)
	input := bufio.NewScanner(c)

	ftp.hello(c)

	for input.Scan() {
		cmd := input.Text()
		log.Println("User input:", cmd)
		ftp.command(c, cmd)
	}
	log.Println("Client disconnected")
}

func loop(listener net.Listener, stopped *bool) {
	for {
		if conn, err := listener.Accept(); err == nil {
			log.Println("Client connected!")
			go handler(conn)
		}
		// stop after at least one round
		if *stopped {
			return
		}
	}
}
