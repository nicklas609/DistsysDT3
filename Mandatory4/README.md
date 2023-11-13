**_How to run program:_**

*Prerequisites*
- Install Consul:
      https://github.com/hashicorp/consul

*Getting Started*
- Start Consul :
    consul agent -dev
- Run Client 1:
    ```go run Client/Client.go "Node 1" :5001 localhost:8500```
- Run Client 2:
    ```go run Client/Client.go "Node 2" :5002 localhost:8500```
- Run Client 3:
    ```go run Client/Client.go "Node 3" :5003 localhost:8500```
- (Run Client x: ```go run Client/Client.go "Node x" :500x localhost:8500``` )

*Request to enter the critical section*
- Type ```critical``` and hit enter 
  
