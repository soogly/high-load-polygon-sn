package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/gobwas/httphead"
	"github.com/gobwas/ws"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	var (
		expectHost = "github.com"
		expectURI  = "/websocket"
	)

	var id int
	reqID := []string{"0"}
	header := http.Header{
		"X-Request-ID": reqID,
	}

	u := ws.ConnUpgrader{
		OnRequest: func(host, uri []byte) (err error, code int) {
			if !bytes.Equal(host, expectHost) {
				return fmt.Errorf("unexpected host: %s", host), 403
			}
			if !bytes.Equal(uri, expectURI) {
				return fmt.Errorf("unexpected uri: %s", uri), 403
			}
			return // Continue upgrade.
		},
		OnHeader: func(key, value []byte) (err error, code int) {
			if !bytes.Equal(key, headerCookie) {
				return
			}
			cookieOK := httphead.ScanCookie(value, func(key, value []byte) bool {
				// Check session here or do some other stuff with cookies.
				// Maybe copy some values for future use.
			})
			if !cookieOK {
				return fmt.Errorf("bad cookie"), 400
			}
			return
		},
		BeforeUpgrade: func() (headerWriter func(io.Writer), err error, code int) {
			// Final checks here before return 101 Continue.

			reqID[0], err = strconv.FormatInt(id, 10)
			if err != nil {
				return nil, err, 500
			}

			return func(w io.Writer) {
				header.Write(w)
			}, nil, 0
		},
	}

	for ; ; id++ {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		_, err := u.Upgrade(conn)
		if err != nil {
			log.Printf("upgrade error: %s", err)
		}
	}
}
