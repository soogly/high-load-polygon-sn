package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"runtime"

	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
)

type Packet struct {
	mess string
}

// Channel wraps user connection.
type Channel struct {
	conn net.Conn    // WebSocket connection.
	send chan Packet // Outgoing packets queue.
}

var connBank []*net.Conn

func main() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		// handle error
	}

	// Prepare handshake header writer from http.Header mapping.
	header := ws.HandshakeHeaderHTTP(http.Header{
		"X-Go-Version": []string{runtime.Version()},
	})

	u := ws.Upgrader{
		OnHost: func(host []byte) error {
			if string(host) == "localhost:8080" {
				println(string(host), "<= host")
				return nil
			}
			println(string(host), "<= ok host")

			return ws.RejectConnectionError(
				ws.RejectionStatus(403),
				ws.RejectionHeader(ws.HandshakeHeaderString(
					"X-Want-Host: localhost\r\n",
				)),
			)
		},
		OnHeader: func(key, value []byte) error {
			if string(key) != "Cookie" {
				return nil
			}
			ok := httphead.ScanCookie(value, func(key, value []byte) bool {
				// Check session here or do some other stuff with cookies.
				// Maybe copy some values for future use.
				return true
			})
			if ok {
				return nil
			}
			return ws.RejectConnectionError(
				ws.RejectionReason("bad cookie"),
				ws.RejectionStatus(400),
			)
		},
		OnBeforeUpgrade: func() (ws.HandshakeHeader, error) {
			return header, nil
		},
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		_, err = u.Upgrade(conn)
		if err != nil {
			log.Printf("upgrade error: %s", err)
		}
		connBank = append(connBank, &conn)
		go readd()
	}

}

var buf *bufio.Reader

func readd() {
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Print("EOF")
				break
			} else {
				fmt.Println(err)
				return
			}
		}
		fmt.Print(line)
	}
}
