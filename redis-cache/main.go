package main

import (
	"redis-lite/command"
	server2 "redis-lite/server"
)

func main() {
	server := server2.Server{Operation: command.Operation{}}
	server.Start()
}