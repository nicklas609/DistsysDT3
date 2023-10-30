package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	proto "github.com/nicklas609/DistsysDT3.git/tree/main/Mandatory3/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	id         int
	portNumber int
}

var (
	clientPort      = flag.Int("cPort", 0, "client port number")
	serverPort      = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
	ClientTimeStamp = int64(1)
)

var connected = false

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
	serverConnection, _ := connectToServer()
	connected = true

	stream, err := serverConnection.PublishReceive(context.Background())

	// Send I am connected message to server
	connectedMessage(client, stream)

	if err != nil {
		print("what")
	}
	go sendMessage(client, stream)
	go publishMessage(client, stream)

	for connected {

	}

	// Send I have left message to server
	disconnectedMessage(client, stream)

	time.Sleep(1 * time.Second)

}

func sendMessage(client *Client, stream proto.Broadcast_PublishReceiveClient) {

	//stream, err := serverConnection.PublishReceive(context.Background())

	// if err != nil {

	// 	print("what")
	// }
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		// publishReturnMessage, err := serverConnection.PublishReceive(context.Background(), &proto.Publish{
		// 	ClientId: int64(client.id),
		// 	Content:  input,
		// })
		input := scanner.Text()
		ClientTimeStamp++

		if input == "!quit" {
			connected = false
		} else {
			message := &proto.Publish{

				ClientId:  int64(client.id),
				Content:   input,
				TimeStamp: ClientTimeStamp,
			}

			//var message = serverConnection.PublishReceive(context).Publish

			if err := stream.Send(message); err != nil {
				log.Fatalf("Failed to send a note: %v", err)
			}
		}

	}

}

func publishMessage(client *Client, stream proto.Broadcast_PublishReceiveClient) {

	// stream, err := serverConnection.PublishReceive(context.Background())
	// if err != nil {

	// 	print("what")
	// }

	for {
		//log.Printf("I am here")
		in, err := stream.Recv()
		// if err == io.EOF {
		// 	// read done.
		// 	close(waitc)
		// 	return
		// }
		if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
		}

		if in.TimeStamp > ClientTimeStamp {

			ClientTimeStamp = in.TimeStamp
		}
		log.Print("Participant ", in.ClientId, " ", in.Content, " at Lamport time ", in.TimeStamp)
	}

}

func connectedMessage(client *Client, stream proto.Broadcast_PublishReceiveClient) {
	message := &proto.Publish{

		ClientId:  int64(client.id),
		Content:   "joined Chitty-Chat",
		TimeStamp: int64(0),
	}

	//var message = serverConnection.PublishReceive(context).Publish

	if err := stream.Send(message); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}
}

func disconnectedMessage(client *Client, stream proto.Broadcast_PublishReceiveClient) {
	message := &proto.Publish{

		ClientId:  int64(client.id),
		Content:   "left Chitty-Chat",
		TimeStamp: ClientTimeStamp,
	}

	//var message = serverConnection.PublishReceive(context).Publish

	if err := stream.Send(message); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}

	// waitc <- 1
}

// func publishMessage(client *Client) {
// 	serverConnection, _ := connectToServer()

// 	scanner := bufio.NewScanner(os.Stdin)

// 	for scanner.Scan() {

// 		input := scanner.Text()

// 		log.Printf("Client wants to publish a message: %s\n", input)

// 		publishReturnMessage, err := serverConnection.PublishReceive(context.Background(), &proto.Publish{
// 			ClientId: int64(client.id),
// 			Content:  input,
// 		})

// 		if err != nil {
// 			log.Printf(err.Error())
// 		} else {
// 			log.Printf("Server %s says the message is %s\n", publishReturnMessage.ServerName, publishReturnMessage.Content)
// 		}
// 	}
// }

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
