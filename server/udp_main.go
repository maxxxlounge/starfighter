package main

import (
    "fmt"
    "net"
	"time"
)

func main() {
    // Resolve UDP address (host:port)
    addr, err := net.ResolveUDPAddr("udp", "localhost:8889")
    if err != nil {
        fmt.Println("Error resolving UDP address:", err)
        return
    }

    // Create UDP connection
    conn, err := net.ListenUDP("udp", addr)
    if err != nil {
        fmt.Println("Error creating UDP connection:", err)
        return
    }
    defer conn.Close()

    fmt.Println("UDP server started. Listening on", addr)

    // Buffer to read incoming data
    buffer := make([]byte, 1024)

    for {
        
		// Read data from the connection
        n, clientAddr, err := conn.ReadFromUDP(buffer)
        if err != nil {
            fmt.Println("Error reading data:", err)
            continue
        }

        // Convert received data to string
        data := string(buffer[:n])
        fmt.Printf("Received from %s: %s\n", clientAddr, data)

		// Respond back to the client (echo)
		// create datetime hello world
		ts := time.Now().Format("2006-01-02 15:04:05")
		
		data = ts +  "Hello from server"
        _, err = conn.WriteToUDP([]byte("Server: "+data), clientAddr)
        if err != nil {
            fmt.Println("Error sending response:", err)
            continue
        }

		time.Sleep(1 * time.Second)
	}
}
