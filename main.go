package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	// Set up a TCP listener on port 8080
	const HOST string = "localhost"
	const PORT int = 8080
	clients := make(map[net.Conn]string)
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", HOST, PORT))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("Listening on localhost:8080")
	for {
		// Accept incoming connections from clients
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			return
		}
		var clientPort string = getPortFromAddress(conn.RemoteAddr().String())
		fmt.Println("Received connection from", clientPort)

		go handleConnection(conn, &clients)
	}
}

func getPortFromAddress(address string) string {
	addressSlice := strings.Split(address, ":")
	var port string = addressSlice[len(addressSlice)-1]
	return port
}

func handleConnection(conn net.Conn, clients *map[net.Conn]string) {
	buffer := make([]byte, 1028)
	var clientPort string = getPortFromAddress(conn.RemoteAddr().String())
	for {
		dataLength, err := conn.Read(buffer)
		if err != nil {
			delete(*clients, conn)
			conn.Close()
			log.Printf("Client %v disconnected\n", clientPort)
			break
		}
		dataBytes := buffer[:dataLength]
		data := string(dataBytes)
		// first message received (register user)
		userName, clientExists := (*clients)[conn]
		if !clientExists {
			(*clients)[conn] = data
			continue
		}

		// end connection
		if data == "END" {
			delete(*clients, conn)
			conn.Close()
			log.Println("Client", clientPort, "disconnected")
			break
		}

		// message received
		message := fmt.Sprintf("[%v] %v\n", userName, data)
		log.Print(message)
		for client := range *clients {
			if conn == client {
				continue
			}
			_, err = client.Write([]byte(message))
			if err != nil {
				log.Print("Error sending data to socket", err)
				conn.Close()
				break
			}
		}
	}
}
