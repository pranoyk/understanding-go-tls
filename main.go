package main

import (
	"bytes"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	// username := os.Getenv("USERNAME")
	// password := os.Getenv("PASSWORD")

	// set net connecction
	netConn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), time.Second*15)
	if err != nil {
		log.Fatalf("error occurred while setting up net connection, error: %v", err)
	}
	defer netConn.Close()

	tlsConfig := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: true,
	}

	// setup client side connection and it returns a greeting which is parsed in the next step
	conn := tls.Client(netConn, tlsConfig)
	defer conn.Close()

	// read greeting
	var buf bytes.Buffer
	var header int32

	err = binary.Read(conn, binary.BigEndian, &header)
	if err != nil {
		log.Fatalf("error while reading bytes from the connection, error: %v", err)
	}

	if header < 4 {
		log.Fatal("error response length is less than 4 which is part of the header")
	}

	msgLen := int64(header - 4)

	// reads the bytes from conn upto the msgLen. Returns when EOF is returned or msgLen is reached
	lr := io.LimitReader(conn, msgLen)

	// reads from lr to buf
	_, err = buf.ReadFrom(lr)
	if err != nil {
		log.Fatalf("error while reading into the buf from io.Reader, error: %v", err)
	}

	fmt.Printf("data: %s", buf.String())
}
