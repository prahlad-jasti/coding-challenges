package main

import (
    "fmt"
    "net"
)

func main() {
    // Define the address to connect to
    address := "localhost:6379"

    // Connect to the server
    conn, err := net.Dial("tcp", address)
    if err != nil {
        fmt.Println("Error connecting:", err)
        return
    }
    defer conn.Close()

    // Send input to the server
    input := "-1\r\n"
    _, err = fmt.Fprintf(conn, input)
    if err != nil {
        fmt.Println("Error sending data:", err)
        return
    }

    fmt.Println("Data sent successfully.")
}