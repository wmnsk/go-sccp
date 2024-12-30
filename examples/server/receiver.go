// Copyright 2018-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

/*
Command receiver receives SCCP messages from the client and prints them out.
*/
package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/wmnsk/go-m3ua/messages/params"
	"github.com/wmnsk/go-sccp"

	"github.com/ishidawataru/sctp"
	"github.com/wmnsk/go-m3ua"
)

func serve(conn *m3ua.Conn) {
	defer conn.Close()

	buf := make([]byte, 1500)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Printf("Closed M3UA conn with: %s, waiting to come back...", conn.RemoteAddr())
				return
			}
			log.Printf("Error reading from M3UA conn: %s", err)
			return
		}

		msg, err := sccp.ParseMessage(buf[:n])
		if err != nil {
			log.Printf("Failed to parse SCCP message: %s, %x", err, buf[:n])
			continue
		}

		log.Printf("Received SCCP message: %s", msg)
	}
}

func main() {
	var (
		addr = flag.String("addr", "127.0.0.1:2905", "Source IP and Port listen.")
	)
	flag.Parse()

	// create *Config to be used in M3UA connection
	config := m3ua.NewServerConfig(
		&m3ua.HeartbeatInfo{
			Enabled:  true,
			Interval: 0,
			Timer:    time.Duration(5 * time.Second),
		},
		0x22222222,                  // OriginatingPointCode
		0x11111111,                  // DestinationPointCode
		1,                           // AspIdentifier
		params.TrafficModeLoadshare, // TrafficModeType
		0,                           // NetworkAppearance
		0,                           // CorrelationID
		[]uint32{1, 2},              // RoutingContexts
		params.ServiceIndSCCP,       // ServiceIndicator
		0,                           // NetworkIndicator
		0,                           // MessagePriority
		1,                           // SignalingLinkSelection
	)
	config.AspIdentifier = nil
	config.CorrelationID = nil

	// setup SCTP listener on the specified IPs and Port.
	laddr, err := sctp.ResolveSCTPAddr("sctp", *addr)
	if err != nil {
		log.Fatalf("Failed to resolve SCTP address: %s", err)
	}

	listener, err := m3ua.Listen("m3ua", laddr, config)
	if err != nil {
		log.Fatalf("Failed to listen: %s", err)
	}
	log.Printf("Waiting for connection on: %s", listener.Addr())

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for {
		conn, err := listener.Accept(ctx)
		if err != nil {
			log.Fatalf("Failed to accept M3UA: %s", err)
		}
		log.Printf("Connected with: %s", conn.RemoteAddr())

		go serve(conn)
	}
}
