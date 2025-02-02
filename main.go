package main

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"time"
)

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	addr, err := url.Parse(host)
	if err != nil {
		log.Fatalf("error occurred while parsing host, error: %v", err)
	}
	isPlainTCP := false
	if addr.Scheme == "" {
		isPlainTCP = true
	}

	// set net connecction
	var parsedHost string
	if isPlainTCP {
		parsedHost = host
	} else {
		parsedHost = addr.Host
	}
	netConn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", parsedHost, port), time.Second*15)
	if err != nil {
		log.Fatalf("error occurred while setting up net connection, error: %v", err)
	}
	defer netConn.Close()

	finalConn := netConn
	if !isPlainTCP {
		tlsConfig := &tls.Config{
			ServerName:         host,
			InsecureSkipVerify: true,
		}

		// setup client side connection and it returns a greeting which is parsed in the next step
		finalConn := tls.Client(netConn, tlsConfig)
		defer finalConn.Close()
	}
	// read greeting
	var buf bytes.Buffer
	var header int32

	err = binary.Read(finalConn, binary.BigEndian, &header)
	if err != nil {
		log.Fatalf("error while reading bytes from the connection, error: %v", err)
	}

	if header < 4 {
		log.Fatal("error response length is less than 4 which is part of the header")
	}

	msgLen := int64(header - 4)

	// reads the bytes from conn upto the msgLen. Returns when EOF is returned or msgLen is reached
	lr := io.LimitReader(finalConn, msgLen)

	// reads from lr to buf
	_, err = buf.ReadFrom(lr)
	if err != nil {
		log.Fatalf("error while reading into the buf from io.Reader, error: %v", err)
	}

	fmt.Printf("data: %s", buf.String())
}
