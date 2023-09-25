package main

import (
	"fmt"
	"sort"
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
	Amount_packages := 6

	syn_ack <- seq

	var res_seq = <-syn_ack
	var ack = <-syn_ack

	if res_seq == seq+1 {
		seq += 1
		ack += 1
		syn_ack <- ack
		syn_ack <- seq

		//var ack_array [6]int
		//var datapackage_array [6]datapackage

		var ack_slice = make([]int, Amount_packages)
		var datapackage_slice = make([]datapackage, Amount_packages)

		for i := 0; i < len(datapackage_slice); i++ {
			datapackage_slice[i] = datapackage{21, 53, ack + 1, seq, messages[i], 1, Amount_packages}
			data_chan <- datapackage_slice[i]
			ack = <-syn_ack
			ack_slice[i] = ack
		}

		for i := 0; i < len(datapackage_slice); i++ {
			var ok = false
			for j := 0; j < len(ack_slice); j++ {
				if ack_slice[j] == datapackage_slice[i].seq_num {
					ok = true
					break
				}
			}
			if !ok {
				data_chan <- datapackage_slice[i]
			}
		}
	}
}

func server() {
	packages := map[datapackage]bool{}

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

		//var datapackagearray [6]datapackage
		//datapackagearray[0] = received_data

		var datapackageSlice = make([]datapackage, received_data.length)
		datapackageSlice[0] = received_data

		for i := 1; i < received_data.length; i++ {
			received_data = <-data_chan
			syn_ack <- received_data.seq_num

			if packages[received_data] {
				//if dublicate go back in the interator
				i--

			} else {
				packages[received_data] = true
				datapackageSlice[i] = received_data
			}

		}

		sort.Slice(datapackageSlice, func(i, j int) bool {
			return datapackageSlice[i].seq_num < datapackageSlice[j].seq_num
		})
		for _, v := range datapackageSlice {
			fmt.Println(v.data)
		}

	}
}

func main() {

	go client()
	go server()
	time.Sleep(time.Second * 10)
}
