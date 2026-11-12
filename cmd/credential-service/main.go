package main

import (
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	log.Println("Credential Service listening on :50051")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go func() {
			defer conn.Close()
			// Handle the connection here.
		}()
	}
}
