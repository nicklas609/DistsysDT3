### How to run program: ###

**Prerequisites** 

- Install Consul: https://github.com/hashicorp/consul 

**Getting Started**

- Start Consul: ```consul agent -dev```
- Run Node 1: ```go run Client/Client.go "Node 1" :5001 localhost:8500```
- Run Node 2: ```go run Client/Client.go "Node 2" :5002 localhost:8500```
- (Run Node x: ```go run Client/Client.go "Node x" :500x localhost:8500```)

- Run User 1: ```go run Node/node.go "User 1" :5011 localhost:8500```
- Run User 2: ```go run Node/node.go "User 2" :5012 localhost:8500```
- (Run User x: ```go run Node/node.go "User x" :501x localhost:8500```)

**How to Bid with a running user:** 
- Step 1: insert ```2``` in the console to request a bid
- Step 2: insert the amount of the bid to the console
