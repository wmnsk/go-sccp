// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Command client sends given payload on top of SCCP UDT message.
package main

import (
	"context"
	"encoding/hex"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ishidawataru/sctp"
	m3params "github.com/wmnsk/go-m3ua/messages/params"

	"github.com/wmnsk/go-sccp"
	"github.com/wmnsk/go-sccp/params"
	"github.com/wmnsk/go-sccp/utils"

	"github.com/wmnsk/go-m3ua"
)

func main() {
	var (
		addr  = flag.String("addr", "127.0.0.1:2905", "Remote IP and Port to connect to.")
		data  = flag.String("data", "deadbeef", "Payload to send on UDT in hex format.")
		hbInt = flag.Duration("hb-interval", 0, "Interval for M3UA BEAT. Put 0 to disable")
	)
	flag.Parse()

	// setup data to send
	payload, err := hex.DecodeString(*data)
	if err != nil {
		log.Fatalf("Failed to decode hex string: %s", err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT)

	// create *Config to be used in M3UA connection
	m3config := m3ua.NewConfig(
		0x11111111,              // OriginatingPointCode
		0x22222222,              // DestinationPointCode
		m3params.ServiceIndSCCP, // ServiceIndicator
		0,                       // NetworkIndicator
		0,                       // MessagePriority
		1,                       // SignalingLinkSelection
	)
	m3config. // set parameters to use
			EnableHeartbeat(*hbInt, 10*time.Second).
			SetAspIdentifier(1).
			SetTrafficModeType(m3params.TrafficModeLoadshare).
			SetNetworkAppearance(0).
			SetRoutingContexts(1, 2)

	// setup SCTP peer on the specified IPs and Port.
	raddr, err := sctp.ResolveSCTPAddr("sctp", *addr)
	if err != nil {
		log.Fatalf("Failed to resolve SCTP address: %s", err)
	}

	// setup underlying SCTP/M3UA connection first
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m3conn, err := m3ua.Dial(ctx, "m3ua", nil, raddr, m3config)
	if err != nil {
		log.Fatal(err)
	}

	cdPA, err := utils.StrToSwappedBytes("1234567890123456", "0")
	if err != nil {
		log.Fatal(err)
	}
	cgPA, err := utils.StrToSwappedBytes("9876543210", "0")
	if err != nil {
		log.Fatal(err)
	}

	// create UDT message with CdPA, CgPA and payload
	udt, err := sccp.NewUDT(
		1,    // Protocol Class
		true, // Message handling
		params.NewPartyAddress( // CalledPartyAddress: 1234567890123456
			0x12, 0, 6, 0x00, // Indicator, SPC, SSN, TT
			0x01, 0x01, 0x04, // NP, ES, NAI
			cdPA, // GlobalTitleInformation
		),
		params.NewPartyAddress( // CallingPartyAddress: 9876543210
			0x12, 0, 7, 0x01, // Indicator, SPC, SSN, TT
			0x01, 0x02, 0x04, // NP, ES, NAI
			cgPA, // GlobalTitleInformation
		),
		payload, // payload
	).MarshalBinary()
	if err != nil {
		log.Fatal(err)
	}

	// send once
	if _, err := m3conn.Write(udt); err != nil {
		log.Fatal(err)
	}

	// send once in 3 seconds
	for {
		ticker := time.NewTicker(3 * time.Second)
		select {
		case sig := <-sigCh:
			log.Printf("Got signal: %v, exitting...", sig)
			ticker.Stop()
			os.Exit(1)
		case <-ticker.C:
			if _, err := m3conn.Write(udt); err != nil {
				log.Fatal(err)
			}
		}
	}
}
