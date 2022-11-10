package app

import (
	"log"
	"net"
)

const PacketSize = 65535

func (app *App) receiver() {
	udpAddr := &net.UDPAddr{
		Port: 21000,
		IP:   net.ParseIP("0.0.0.0"),
	}

	ln, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		message := make([]byte, PacketSize)
		n, err := ln.Read(message)
		if err != nil || n == 0 {
			log.Println("error read packet", err)
			continue
		}

		app.message <- message[0:n]
	}
}
