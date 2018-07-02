package main

import (
	"net"
	"log"
	"bufio"
	"fmt"
)

var Users = make(map[net.Conn]string)
var newConnections = make(chan net.Conn)
var messages = make(chan string)

func handleClient(connection net.Conn) {

	fmt.Fprint(connection, "Enter user name: ")

	userName := bufio.NewScanner(connection)
	userName.Scan()

	Users[connection] = userName.Text()

	log.Println(userName.Text() + " connected")

	fmt.Fprintln(connection, "Hi "+userName.Text())

	for {

		buffer, err := bufio.NewReader(connection).ReadString('\n')

		if err != nil {
			log.Println("User " + userName.Text() + " disconnected")
			delete(Users, connection)
			connection.Close()
			return
		}

		messages <- fmt.Sprintf(userName.Text() + ": " + buffer)
		log.Print(userName.Text() + ": " + buffer)
	}

}

func handleMessage(connection net.Conn, message string, userName string) {

	_, err := connection.Write([]byte(message))

	if err != nil {
		log.Println("User " + userName + " disconnected")
		delete(Users, connection)
		connection.Close()
		return
	}
}

func acceptNewClient(server net.Listener) {

	for {
		client, err := server.Accept() //wait for user

		if err != nil {
			log.Println("User can't join to server. Error: ", err.Error())
			return
		}
		newConnections <- client //write client to newConnections canal
	}
}

func main() {

	log.Println("Server is running!")

	server, err := net.Listen("tcp", ":8080")

	if err != nil {
		log.Println("Can't start server! Error: ", err.Error())
	}

	go acceptNewClient(server)

	for {
		select {

		case connection := <-newConnections:
			go handleClient(connection)

		case message := <-messages:
			for client, userName := range Users {
				handleMessage(client, message, userName)
			}
		}
	}
}
