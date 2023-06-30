package server

import (
	"log"
	"net"
)

type TeltonikaServer struct {
	DataChan chan []byte
}

func NewTeltonikaServer() *TeltonikaServer {
	return &TeltonikaServer{
		DataChan: make(chan []byte),
	}
}

func (s *TeltonikaServer) ServeUDP(addr string) {
	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer pc.Close()

	buffer := make([]byte, 1024)

	for {
		n, addr, err := pc.ReadFrom(buffer)
		if err != nil {
			log.Println(err)
			continue
		}

		go s.handleConnection(pc, addr, buffer[:n])
	}
}

func (s *TeltonikaServer) handleConnection(pc net.PacketConn, addr net.Addr, data []byte) {
	s.DataChan <- data
}
