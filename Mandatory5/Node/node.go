package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	proto "github.com/nicklas609/DistsysDT3/tree/main/Mandatory5/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func setLog() *os.File {
	// Clears the log.txt file when a new server is started
	// if err := os.Truncate("log.txt", 0); err != nil {
	// 	log.Printf("Failed to truncate: %v", err)
	// }

	// This connects to the log file/changes the output of the log information to the log.txt file.
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(f)
	return f
}

type Client struct {
	proto.UnimplementedCriticalServiceServer

	// Self information
	Name string
	Addr string

	// Consul related variables
	SDAddress string
	SDKV      api.KV

	// used to make requests
	Clients      map[string]proto.CriticalServiceClient
	Users        map[string]proto.CriticalServiceClient
	NodesReplies map[string]bool
	InCritSys    bool
	timeStamp    int64
	Username     string
}

func (c *Client) GetnodeType(ctx context.Context, in *proto.Ack) (*proto.NodeType, error) {
	return &proto.NodeType{Type: false}, nil
}

func getRes(c *Client) {

	amount := 0
	for key, element := range c.Clients {
		amount++
		if key == "nil" && element == nil {
			log.Print("Go is weird and I don't have time to find a better way for this")
		}
	}

	temp := 0
	RequestNode := rand.Intn((amount) + 1)
	if RequestNode == 0 {
		RequestNode++
	}
	for key, element := range c.Clients {
		temp++
		if temp == RequestNode {

			r, t := element.GetResult(context.Background(), &proto.AskForResult{Res: "What is the current result?"})
			if t != nil || key == "nil" {

			}
			log.Print(r.Result)
		}

	}

}

func makeBid(c *Client, bid int64) {

	amount := 0
	for key, element := range c.Clients {
		amount++
		if key == "nil" && element == nil {
			log.Print("Go is weird and I don't have time to find a better way for this")
		}
	}

	temp := 0
	RequestNode := rand.Intn((amount) + 1)
	if RequestNode == 0 {
		RequestNode++
	}
	for key, element := range c.Clients {
		temp++
		if temp == RequestNode {

			r, t := element.MakeBid(context.Background(), &proto.Bid{Amount: bid, Bidder: c.Username})
			if t != nil || key == "nil" {

			}
			log.Print(r.Message)
		}

	}

}

// Start listening/service.
func (c *Client) StartListening() {

	lis, err := net.Listen("tcp", c.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	_c := grpc.NewServer() // n is for serving purpose

	proto.RegisterCriticalServiceServer(_c, c)
	// Register reflection service on gRPC server.
	reflection.Register(_c)

	// start listening
	if err := _c.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Register self with the service discovery module.
// This implementation simply uses the key-value store. One major drawback is that when nodes crash. nothing is updated on the key-value store. Services are a better fit and should be used eventually.
func (c *Client) registerService() {
	config := api.DefaultConfig()
	config.Address = c.SDAddress
	consul, err := api.NewClient(config)
	if err != nil {
		log.Panicln("Unable to contact Service Discovery.")
	}

	kv := consul.KV()
	p := &api.KVPair{Key: c.Name, Value: []byte(c.Addr)}
	_, err = kv.Put(p, nil)
	if err != nil {
		log.Panicln("Unable to register with Service Discovery.")
	}

	// store the kv for future use
	c.SDKV = *kv

	log.Println("Successfully registered with Consul.")
}

// Start the node.
// This starts listening at the configured address. It also sets up clients for it's peers.
func (c *Client) Start() {
	// init required variables
	c.Clients = make(map[string]proto.CriticalServiceClient)
	c.Users = make(map[string]proto.CriticalServiceClient)
	c.NodesReplies = make(map[string]bool)
	c.InCritSys = false
	c.timeStamp = 1
	c.Username = ""
	//f := setLog()

	// start service / listening
	go c.StartListening()

	// register with the service discovery unit
	c.registerService()

	c.GreetAll()

	//log.Print("I am here too")

	go menu(c)

	for {
		time.Sleep(3 * time.Second)
		c.GreetAll()

		// for key, element := range c.Clients {
		// 	if key != "Why do I need to use key!!!!!" {
		// 		r, _ := element.RequestCritical(context.Background(), &proto.Request{Name: "Hello"})
		// 		log.Print(r.Message)
		// 	}
		// }

	}
	//defer f.Close()
}

func menu(c *Client) {

	scanner := bufio.NewScanner(os.Stdin)

	log.Print("What is your name")
	scanner.Scan()
	name := scanner.Text()
	c.Username = name

	log.Printf("Type 1 to see current bid")
	log.Printf("Type 2 to make a bid")
	for scanner.Scan() {

		input := scanner.Text()

		if input == "1" {

			getRes(c)

		} else if input == "2" {
			log.Print("How much do you want to bid?")
			scanner.Scan()
			bid := scanner.Text()
			i, err := strconv.Atoi(bid)

			if err != nil {
				log.Print("Please input a number")
			}

			makeBid(c, int64(i))
		}

	}
}

// Setup a new grpc client for contacting the server at addr.
func (n *Client) SetupClient(name string, addr string) {

	// setup connection with other node
	log.Print("What is going off2")
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	//defer conn.Close()
	r, t := proto.NewCriticalServiceClient(conn).GetnodeType(context.Background(), &proto.Ack{Message: "Whats your type"})

	if t != nil {
		log.Fatalf("Something is wrong here")
	}

	if r.Type {
		n.Clients[name] = proto.NewCriticalServiceClient(conn)
		n.NodesReplies[name] = false
		n.timeStamp++
	} else {

		n.Users[name] = proto.NewCriticalServiceClient(conn)
	}

}

// Busy Work module, greet every new member you find
func (c *Client) GreetAll() {
	// get all nodes -- inefficient, but this is just an example
	kvpairs, _, err := c.SDKV.List("Node", nil)
	if err != nil {
		log.Panicln(err)
	}

	// fmt.Println("Found nodes: ")
	for _, kventry := range kvpairs {
		if strings.Compare(kventry.Key, c.Name) == 0 {
			// ourself
			continue
		}
		if c.Clients[kventry.Key] == nil && c.Users[kventry.Key] == nil {
			fmt.Println("New member: ", kventry.Key)
			// connection not established previously
			c.SetupClient(kventry.Key, string(kventry.Value))
		}
	}
}

func main() {
	// pass the port as an argument and also the port of the other node
	args := os.Args[1:]

	if len(args) < 3 {
		fmt.Println("Arguments required: <name> <listening address> <consul address>")
		os.Exit(1)
	}

	// args in order
	name := args[0]
	listenaddr := args[1]
	sdaddress := args[2]

	noden := Client{Name: name, Addr: listenaddr, SDAddress: sdaddress, Clients: nil} // noden is for opeartional purposes

	// start the node
	noden.Start()
}
