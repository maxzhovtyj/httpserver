package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}

	udp, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	defer udp.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		readString, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		_, err = udp.Write([]byte(readString))
		if err != nil {
			log.Fatal(err)
		}
	}
}
