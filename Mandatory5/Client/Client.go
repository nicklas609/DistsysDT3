package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	proto "github.com/nicklas609/DistsysDT3/tree/main/Mandatory5/proto"

	"golang.org/x/net/context"
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
	Clients       map[string]proto.CriticalServiceClient
	NodesReplies  map[string]bool
	InCritSys     bool
	timeStamp     int64
	Leader        string
	viceLeader    string
	IamLeader     bool
	IamViceLeader bool
	MaxBid        int64
	CurrentBidde  string
}

func (c *Client) AreYouTheLeader(ctx context.Context, in *proto.Request) (*proto.Reply, error) {

	return &proto.Reply{Message: "Yes you may " + c.Name, TimeStamp: c.timeStamp, Leader: c.IamLeader}, nil
}

func (c *Client) leaderWrite(ctx context.Context, in *proto.Request) (*proto.ack, error) {

	c.MaxBid = in.Amount
	c.CurrentBidde = in.Bidder

	return &proto.ack{Message: "ack"}, nil
}

func (c *Client) MakeBid(ctx context.Context, in *proto.bid) (*proto.ack, error) {

	if in.amount > c.MaxBid {

		if c.IamLeader == true {
			c.MaxBid = in.amount
			c.CurrentBidde = in.bidder
			writeToNodes(c)

		} else {
			c.Clients[c.Leader].MakeBid(context.Background(), &proto.bid{Amount: in.amount, Bidder: in.bidder})
			// if ack == "ack" {

			// }
		}

	}
	return &proto.ack{Message: "ack"}, nil
}

func writeToNodes(c *Client) {
	for key, element := range c.Clients {
		if key != "Why do I need to use key!!!!!" {
			element.leaderWrite(context.Background(), &proto.bid{Amount: c.MaxBid, Bidder: c.CurrentBidde})
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
	c.NodesReplies = make(map[string]bool)
	c.timeStamp = 1
	c.Leader = ""
	c.viceLeader = ""
	c.IamLeader = false
	c.IamViceLeader = false
	c.MaxBid = 0
	c.CurrentBidde = ""
	//f := setLog()

	// start service / listening
	go c.StartListening()

	// register with the service discovery unit
	c.registerService()
	c.GreetAll()

	if c.IamLeader == false {
		if c.Leader == "" {
			Findleader(c)
		}
	}

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

// func SendMessage(c *Client) {

// 	scanner := bufio.NewScanner(os.Stdin)

// 	for scanner.Scan() {

// 		input := scanner.Text()

// 		for key, element := range c.Clients {
// 			if key != "Why do I need to use key!!!!!" {
// 				r, t := element.RequestCritical(context.Background(), &proto.Request{Name: input})
// 				log.Print(r.Message)
// 				log.Print(t)
// 			}
// 		}

// 	}

// }

// Setup a new grpc client for contacting the server at addr.
func (n *Client) SetupClient(name string, addr string) {

	// setup connection with other node
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	//defer conn.Close()
	n.Clients[name] = proto.NewCriticalServiceClient(conn)
	n.NodesReplies[name] = false
	n.timeStamp++

	//r, err := n.Clients[name].RequestCritical(context.Background(), &proto.Request{Name: n.Name})
	// if err != nil {
	// 	log.Fatalf("could not greet: %v", err)
	// }
	//log.Printf("Greeting from the other node: %s", r.Message)

}

// Busy Work module, greet every new member you find
func (c *Client) GreetAll() {
	// get all nodes -- inefficient
	kvpairs, _, err := c.SDKV.List("Node", nil)
	if err != nil {
		log.Panicln(err)
	}
	amount := 0

	for _, kventry := range kvpairs {
		amount++
		if strings.Compare(kventry.Key, c.Name) == 0 {
			// ourself
			continue
		}
		if c.Clients[kventry.Key] == nil {
			fmt.Println("New member: ", kventry.Key)
			// connection not established previously
			c.SetupClient(kventry.Key, string(kventry.Value))

			if c.IamLeader == true {
				if c.viceLeader == "" {
					c.viceLeader = kventry.Key
				}
			}
		}
	}
	//You are the only node on the network
	if amount == 1 {
		c.IamLeader = true
	}

}

func Findleader(c *Client) {

	for key, element := range c.Clients {
		if key != "Why do I need to use key!!!!!" {
			r, t := element.AreYouTheLeader(context.Background(), &proto.Request{Name: "Are you the leader"})

			log.Print(r.Leader)

			if r == nil {

			}
			if t == nil {

			}

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
