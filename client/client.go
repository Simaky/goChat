package main

import "net"
import (
	"fmt"
	"bufio"
	"os"
	"flag"
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

	var ip = flag.String("ip","localhost","Server IP Address")
	var port = flag.String("port","8080","Server Port")
	flag.Parse()

	connection, err := net.Dial("tcp", *ip + ":" + *port)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		return
	}

	go getMessage(connection)

	sendMessage(connection)
}
