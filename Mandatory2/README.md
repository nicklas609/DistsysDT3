# a - What are packages in your implementation? What data structure do you use to transmit data and meta-data?
* In our implementation we have created our own struct which represents the packages and contains the TCP header and body as field data.

# b - Does your implementation use threads or processes? Why is it not realistic to use threads?
* We use threads in our implementation. This is not realistic, because we will not encounter any errors. One of the reasons being that we use channels to communicate between our go threads.

# c - In case the network changes the order in which messages are delivered, how would you handle message re-ordering?
* To handle possible changes in the order of sent messages, we use a simple sort algorithm to re-order the messages after receiving.
* IMPROVEMENT: Instead of sorting at the end, insert icoming data into a minheap priority queue, where key is sequence. 

# d - In case messages can be delayed or lost, how does your implementation handle message loss?
* To handle message loss, our implementation check if ackowledgements to all messages have been received, and if not the message is re-send.

# e - Why is the 3-way handshake important?
* The 3-way handshake is important, to ensure there is a mutual/agreed connection between server and client.
