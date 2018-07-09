package main

import (
	"net"
	"log"
	"bufio"
	"fmt"
	"sync"
	"flag"
)

type Users struct {
	sync.Mutex
	data map[net.Conn]user
}

type user struct {
	userName string
}

var newConnections = make(chan net.Conn)
var messages = make(chan string)

func handleClient(connection net.Conn, users *Users) {

	var newUser user

	fmt.Fprint(connection, "Enter user name: ")

	userName := bufio.NewScanner(connection)
	userName.Scan()

	newUser.userName = userName.Text()

	users.Lock()
	users.data[connection] = newUser
	users.Unlock()

	log.Println(userName.Text() + " connected")
	messages <- fmt.Sprintln("User " + userName.Text() + " connected")


	for {

		buffer, err := bufio.NewReader(connection).ReadString('\n')

		if err != nil {
			log.Println("User " + userName.Text() + " disconnected")
			messages <- fmt.Sprintln("User " + userName.Text() + " disconnected")
			users.Lock()
			delete(users.data, connection)
			users.Unlock()
			connection.Close()
			return
		}

		messages <- fmt.Sprintf(userName.Text() + ": " + buffer)
		log.Print(userName.Text() + ": " + buffer)

	}

}

func handleMessage(connection net.Conn, message string, userName string, users *Users) {

	_, err := connection.Write([]byte(message))

	if err != nil {
		log.Println("User " + userName + " disconnected")
		messages <- fmt.Sprintln("User " + userName + " disconnected")
		users.Lock()
		delete(users.data, connection)
		users.Unlock()
		connection.Close()
	}
}

func acceptNewClient(server net.Listener) {

	for {
		client, err := server.Accept() //wait for user

		if err != nil {
			log.Println("User can't join to server. Error: ", err.Error())
		}
		newConnections <- client //write client to newConnections canal
	}
}

func main() {

	log.Println("Server is running!")
	var users Users
	users.data = make(map[net.Conn]user)

	var ip = flag.String("ip","localhost","Server IP Address")
	var port = flag.String("port","8080","Server Port")
	flag.Parse()

	server, err := net.Listen("tcp", *ip + ":" + *port)

	if err != nil {
		log.Println("Can't start server! Error: ", err.Error())
	}

	go acceptNewClient(server)

	for {
		select {

		case connection := <-newConnections:
			go handleClient(connection, &users)

		case message := <-messages:
			for client, user := range users.data {
				handleMessage(client, message, user.userName, &users)
			}
		}
	}
}
