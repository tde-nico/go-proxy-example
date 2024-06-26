package main

import (
	"flag"
	"io"
	"log"
	"net"
)

func handleConn(lconn net.Conn, addr *net.TCPAddr) {
	defer lconn.Close()

	rconn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		log.Printf("DialTCP: %v", err)
		return
	}
	defer rconn.Close()

	go func() {
		io.Copy(lconn, rconn)
		log.Printf("Connection to %s closed\n", addr.String())
	}()
	io.Copy(rconn, lconn)

	log.Printf("Connection from %s closed\n", lconn.RemoteAddr())
}

func main() {
	const localAddr = "0.0.0.0"
	var (
		addr string
		from string
		to   string
		help bool
		h    bool
	)

	flag.StringVar(&addr, "addr", "127.0.0.1", "service address")
	flag.StringVar(&from, "from", "5078", "from port")
	flag.StringVar(&to, "to", "22", "to port")
	flag.BoolVar(&help, "help", false, "help")
	flag.BoolVar(&h, "h", false, "help")
	flag.Parse()

	if help || h {
		flag.PrintDefaults()
		return
	}

	listenAddr, err := net.ResolveTCPAddr("tcp", localAddr+":"+from)
	if err != nil {
		log.Fatalf("ResolveTCPAddr: %v", err)
	}
	connectAddr, err := net.ResolveTCPAddr("tcp", addr+":"+to)
	if err != nil {
		log.Fatalf("ResolveTCPAddr: %v", err)
	}

	ln, err := net.ListenTCP("tcp", listenAddr)
	if err != nil {
		log.Fatalf("ListenTCP: %v", err)
	}
	defer ln.Close()

	log.Printf("Listening on %s:%s\n", localAddr, from)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Accept: %v", err)
			continue
		}

		log.Printf("Connection from %s\n", conn.RemoteAddr())
		go handleConn(conn, connectAddr)
	}
}
