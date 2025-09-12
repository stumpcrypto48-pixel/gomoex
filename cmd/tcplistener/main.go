package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	// f, err := os.Open("messages.txt")
	// if err != nil {
	// 	fmt.Printf("Error : %v", err)
	// 	return
	// }

	// // defer f.Close()
	// for chanLine := range getLinesChannel(f) {
	// 	fmt.Printf("read: %s \n ", chanLine)
	// }

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Printf("Error :: %v", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error :: %v", err)
			break
		}
		defer conn.Close()
		fmt.Println("Connection established!")
		linesFromConn := getLinesChannel(conn)
		for line := range linesFromConn {
			fmt.Printf("read: %s \n", line)
		}
		checkChanConnect(linesFromConn)
	}

}

func checkChanConnect(linesFromConn <-chan string) {
	_, connCheck := <-linesFromConn
	if connCheck == false {
		fmt.Println("Connection is over")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {

	lines := make(chan string)
	go func() {
		defer close(lines)
		defer f.Close()

		readBuff := make([]byte, 8)
		currentLine := ""

		for {
			n, err := f.Read(readBuff)
			if err != nil && err != io.EOF {
				fmt.Printf("Error :: %v", err)
				return
			}
			splitedMessage := strings.Split(string(readBuff[:n]), "\n")

			if len(splitedMessage) > 1 {
				currentLine += splitedMessage[0]
				lines <- currentLine
				currentLine = ""
				for _, line := range splitedMessage[1:] {
					currentLine += line
				}
			} else {
				currentLine += string(readBuff[:n])
			}

			if err == io.EOF {
				if currentLine != "" {
					lines <- currentLine
				}
				break
			}
		}

	}()

	return lines
}
