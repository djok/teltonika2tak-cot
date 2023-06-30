package server

import (
	"io/ioutil"
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

func (s *TeltonikaServer) ServeTCP(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *TeltonikaServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	data, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Println(err)
		return
	}

	s.DataChan <- data
}
