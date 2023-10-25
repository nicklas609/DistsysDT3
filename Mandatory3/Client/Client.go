package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"strconv"

	proto "github.com/nicklas609/DistsysDT3.git/tree/main/Mandatory3/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	id         int
	portNumber int
}

var (
	clientPort = flag.Int("cPort", 0, "client port number")
	serverPort = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
)

func main() {
	// Parse the flags to get the port for the client
	flag.Parse()

	// Create a client
	client := &Client{
		id:         1,
		portNumber: *clientPort,
	}

	// Wait for the client (user) to ask for the time
	// go waitForTimeRequest(client)
	go publishMessage(client)

	for {

	}
}

func publishMessage(client *Client) {
	serverConnection, _ := connectToServer()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		input := scanner.Text()

		log.Printf("Client wants to publish a message: %s\n", input)

		publishReturnMessage, err := serverConnection.PublishReceive(context.Background(), &proto.Publish{
			ClientId: int64(client.id),
			Content:  input,
		})

		if err != nil {
			log.Printf(err.Error())
		} else {
			log.Printf("Server %s says the message is %s\n", publishReturnMessage.ServerName, publishReturnMessage.Content)
		}
	}
}

// func waitForTimeRequest(client *Client) {
// 	// Connect to the server
// 	serverConnection, _ := connectToServer()

// 	// Wait for input in the client terminal
// 	scanner := bufio.NewScanner(os.Stdin)
// 	for scanner.Scan() {
// 		input := scanner.Text()
// 		log.Printf("Client asked for time with input: %s\n", input)

// 		// Ask the server for the time
// 		timeReturnMessage, err := serverConnection.PublishReceive(context.Background(), &proto.AskForTimeMessage{
// 			ClientId: int64(client.id),
// 		})

// 		if err != nil {
// 			log.Printf(err.Error())
// 		} else {
// 			log.Printf("Server %s says the time is %s\n", timeReturnMessage.ServerName, timeReturnMessage.Time)
// 		}
// 	}
// }

func connectToServer() (proto.BroadcastClient, error) {
	// Dial the server at the specified port.
	conn, err := grpc.Dial("localhost:"+strconv.Itoa(*serverPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Could not connect to port %d", *serverPort)
	} else {
		log.Printf("Connected to the server at port %d\n", *serverPort)
	}
	return proto.NewBroadcastClient(conn), nil
}
