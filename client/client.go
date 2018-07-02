package main

import "net"
import (
	"fmt"
	"bufio"
	"os"
)

func getMessage(connection net.Conn) {

	for {
		buf := make([]byte, 1024)
		lens, err := connection.Read(buf)
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}

		fmt.Print(string(buf[:lens]))
	}
}

func sendMessage(connection net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')

		fmt.Fprintf(connection, text) //send to conn
	}
}

func main() {

	connection, err := net.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	go getMessage(connection)

	sendMessage(connection)
}
