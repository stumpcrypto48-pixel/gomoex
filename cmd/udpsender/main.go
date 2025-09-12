package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpAddress, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddress)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer udpConn.Close()

	buffReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		str, err := buffReader.ReadString('\n')
		if err != nil {
			break
		}
		udpConn.Write([]byte(str))
	}

}
