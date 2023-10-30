package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	proto "github.com/nicklas609/DistsysDT3.git/tree/main/Mandatory3/proto"
	"google.golang.org/grpc"
)

// Struct that will be used to represent the Server.
type Server struct {
	proto.UnimplementedBroadcastServer // Necessary
	name                               string
	port                               int
}

var ServerTimestamp = int64(1)
var streams = make([]proto.Broadcast_PublishReceiveServer, 0)

// Used to get the user-defined port for the server from the command line
var port = flag.Int("port", 0, "server port number")

func main() {
	// Get the port from the command line when the server is run
	flag.Parse()

	// Create a server struct
	server := &Server{
		name: "serverName",
		port: *port,
	}

	// Start the server
	go startServer(server)

	// Keep the server running until it is manually quit
	for {

	}
}
func setLog() *os.File {
	// Clears the log.txt file when a new server is started
	if err := os.Truncate("log.txt", 0); err != nil {
		log.Printf("Failed to truncate: %v", err)
	}

	// This connects to the log file/changes the output of the log information to the log.txt file.
	f, err := os.OpenFile("log.txt", os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}

func startServer(server *Server) {

	f := setLog()

	// Create a new grpc server
	grpcServer := grpc.NewServer()

	// Make the server listen at the given port (convert int port to string)
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(server.port))

	if err != nil {
		log.Fatalf("Could not create the server %v", err)
	}
	log.Printf("Started server at port: %d\n", server.port)

	// Register the grpc server and serve its listener
	proto.RegisterBroadcastServer(grpcServer, server)
	serveError := grpcServer.Serve(listener)
	if serveError != nil {
		log.Fatalf("Could not serve listener")
	}

	defer f.Close()

}

func (c *Server) AskForTime(ctx context.Context, in *proto.AskForTimeMessage) (*proto.TimeMessage, error) {
	log.Printf("Client with ID %d asked for the time\n", in.ClientId)
	return &proto.TimeMessage{
		Time:       time.Now().String(),
		ServerName: c.name,
	}, nil
}

func (s *Server) PublishReceive(stream proto.Broadcast_PublishReceiveServer) error {

	streams = append(streams, stream)

	for {

		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if in != nil {

			if ServerTimestamp < in.TimeStamp {
				ServerTimestamp = in.TimeStamp
			}
			if in.TimeStamp == 0 {
				in.TimeStamp = ServerTimestamp + 1
			}

			log.Print("ServersideLog : ", "Participant ", in.Clientname, " ", in.Content, " at Lamport time ", in.TimeStamp)
		}

		for _, s := range streams {

			if in.TimeStamp == 0 {
				in.TimeStamp = ServerTimestamp
			}

			s.Send(in)
		}
	}
}
