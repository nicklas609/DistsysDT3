package main

import (
	"fmt"
	"time"
)

var messages = [6]string{"Hey", "please", "recieve", "all", "6", "datapackages"}

type datapackage struct {
	source_port      int
	destination_port int
	seq_num          int
	ack_num          int
	data             string
	checksum         int
	length           int
}

var syn_ack = make(chan int)
var data_chan = make(chan datapackage)

func client() {
	var seq = 0

	syn_ack <- seq

	var res_seq = <-syn_ack
	var ack = <-syn_ack

	if res_seq == seq+1 {
		seq += 1
		ack += 1
		syn_ack <- ack
		syn_ack <- seq

		var ack_array [6]int
		var datapackage_array [6]datapackage

		for i := 0; i < 6; i++ {
			datapackage_array[i] = datapackage{21, 53, ack + 1, seq, messages[i], 1, 6}
			data_chan <- datapackage_array[i]
			ack = <-syn_ack
			ack_array[i] = ack
		}

		for i := 0; i < len(datapackage_array); i++ {
			var ok = false
			for j := 0; j < len(ack_array); j++ {
				if ack_array[j] == datapackage_array[i].seq_num {
					ok = true
					break
				}
			}
			if !ok {
				data_chan <- datapackage_array[i]
			}
		}
	}
}

func server() {
	var ack = <-syn_ack + 1
	var seq = 100

	syn_ack <- ack
	syn_ack <- seq

	var res_seq = <-syn_ack
	var res_ack = <-syn_ack

	if res_seq == seq+1 && res_ack == ack {
		seq += 1
		fmt.Println("Connection Established")

		var received_data = <-data_chan
		syn_ack <- received_data.seq_num

		var datapackagearray [6]datapackage
		datapackagearray[0] = received_data

		for i := 1; i < received_data.length; i++ {
			received_data = <-data_chan
			syn_ack <- received_data.seq_num
			datapackagearray[i] = received_data
		}

		datapackagearray = sort(datapackagearray)

		for i := 0; i < len(datapackagearray); i++ {
			fmt.Println("Sorted data: ", datapackagearray[i].data)
		}
	}
}

func sort(array [6]datapackage) [6]datapackage {
	for i := 0; i < len(array); i++ {

		for j := i; j > 0; j-- {
			if array[j].seq_num < array[j-1].seq_num {
				var temp1 = array[j]
				var temp2 = array[j-1]
				array[j] = temp2
				array[j-1] = temp1
			}
		}
	}
	return array
}

func main() {

	go client()
	go server()
	time.Sleep(time.Second * 10)
}
