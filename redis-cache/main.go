package main

import (
	"fmt"
	"net"
	"log"
	"io"
)

func handleConnection(conn net.Conn) {
	
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Println("error closing connection", err)
		}
		fmt.Println("connection closed", conn.RemoteAddr())
	}()

	fmt.Println("connection received", conn.RemoteAddr())
	buffer := make([]byte, 1024)
	for {
        n, err := conn.Read(buffer)
        if err != nil {
            if err == io.EOF {
                return
            }
            log.Println("Error reading from connection", err)
            return
        }

        fmt.Printf("Received data from %s: %s\n", conn.RemoteAddr(), buffer[:n])
    }
}

func main() {
	fmt.Println("Starting redis server...")
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleConnection(conn)
	}
}
