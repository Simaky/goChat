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
		text, err := reader.ReadString('\n')

		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}

		fmt.Fprintf(connection, text) //send to conn
	}
}

func main() {

	argsAddress := os.Args[1]
	argsPort := os.Args[2]

	connection, err := net.Dial("tcp", argsAddress + ":" + argsPort)

	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	go getMessage(connection)

	sendMessage(connection)
}
