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
	id         string
	portNumber int
}

func setLog() *os.File {
	// Clears the log.txt file when a new server is started
	// if err := os.Truncate("log.txt", 0); err != nil {
	// 	log.Printf("Failed to truncate: %v", err)
	// }

	// This connects to the log file/changes the output of the log information to the log.txt file.
	f, err := os.OpenFile("log.txt", os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}

var (
	clientPort      = flag.Int("cPort", 0, "client port number")
	serverPort      = flag.Int("sPort", 0, "server port number (should match the port used for the server)")
	ClientTimeStamp = int64(1)
)

var connected = false

func main() {
	// Parse the flags to get the port for the client

	f := setLog()
	flag.Parse()
	println("Enter username:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()

	// Create a client
	client := &Client{
		id:         input,
		portNumber: *clientPort,
	}

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
	defer f.Close()

	time.Sleep(1 * time.Second)

}

func sendMessage(client *Client, stream proto.Broadcast_PublishReceiveClient) {

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {

		input := scanner.Text()
		ClientTimeStamp++

		if input == "!quit" {
			connected = false
		} else {
			message := &proto.Publish{

				Clientname: client.id,
				Content:    input,
				TimeStamp:  ClientTimeStamp,
			}

			if err := stream.Send(message); err != nil {
				log.Fatalf("Failed to send a note: %v", err)
			}
		}

	}

}

func publishMessage(client *Client, stream proto.Broadcast_PublishReceiveClient) {

	for {
		in, err := stream.Recv()

		if err != nil {
			log.Fatalf("Failed to receive a note : %v", err)
		}

		if in.TimeStamp > ClientTimeStamp {

			ClientTimeStamp = in.TimeStamp
		}
		log.Print("ClientsideLog ", "From cliet: ", client.id, " : ", "Participant ", in.Clientname, " ", in.Content, " at Lamport time ", in.TimeStamp)
	}

}

func connectedMessage(client *Client, stream proto.Broadcast_PublishReceiveClient) {
	message := &proto.Publish{

		Clientname: client.id,
		Content:    "joined Chitty-Chat",
		TimeStamp:  int64(0),
	}

	if err := stream.Send(message); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}
}

func disconnectedMessage(client *Client, stream proto.Broadcast_PublishReceiveClient) {
	message := &proto.Publish{

		Clientname: client.id,
		Content:    "left Chitty-Chat",
		TimeStamp:  ClientTimeStamp,
	}

	if err := stream.Send(message); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}
}

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
