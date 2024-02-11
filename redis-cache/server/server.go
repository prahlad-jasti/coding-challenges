package server

import (
	"net"
	"redis-lite/command"
	"redis-lite/util"
)

type Server struct {
	Operation command.Operation
}

func (server *Server) Start() {
	ln, err := net.Listen("tcp", ":6379")

	if err != nil {
		panic("could not Start server")
	}
	for {
		con, err := ln.Accept()
		if err != nil {
			panic("could not establish connection")

		}
		go server.handleConnections(con)
	}
}

func (server *Server) handleConnections(con net.Conn) {
	defer con.Close()
	buffer := make([]byte, 1024)
	_, err := con.Read(buffer)
	if err != nil {
		//fmt.Println("Error reading:", err)
		return
	}

	resp := server.Operation.Handle(buffer)
	con.Write(util.Byte(resp))
}