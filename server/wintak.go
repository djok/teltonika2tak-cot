package server

import (
	"fmt"
	"log"
	"net"
)

type WintakServer struct {
	Address string
}

func NewWintakServer(address string) *WintakServer {
	return &WintakServer{
		Address: address,
	}
}

func (s *WintakServer) Send(data []byte) {
	conn, err := net.Dial("tcp", s.Address)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(data)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Sent data to WINTAK server")
}
