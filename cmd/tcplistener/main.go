package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	tcpSocket, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal(err)
	}

	defer tcpSocket.Close()

	for {
		conn, err := tcpSocket.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("A connection has been accepted.")

		req, err := request.ParseFromReader(conn)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(req.String())

		fmt.Println("The connection has been closed when the channel is closed.")
	}
}
