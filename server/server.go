package server

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gophergala/honeybee/protobee"
	"net"
)

func handleProtoClient(conn net.Conn, c chan *protobee.Connections) {
	fmt.Println("Connection established")
	//Close the connection when the function exits
	defer conn.Close()
	//Create a data buffer of type byte slice with capacity of 4096
	data := make([]byte, 4096)
	//Read the data waiting on the connection and put it in the data buffer
	n, err := conn.Read(data)

	if err != nil {
		fmt.Println("error", err)
	}

	fmt.Println("Decoding Protobuf message")
	//Create an struct pointer of type ProtobufTest.TestMessage struct
	protodata := new(protobee.Connections)
	//Convert all the data retrieved into the ProtobufTest.TestMessage struct type
	err = proto.Unmarshal(data[0:n], protodata)
	if err != nil {
		fmt.Println("error", err)
	}
	//Push the protobuf message into a channel
	c <- protodata
}

func Run() {
	fmt.Println("listen for pb packets")
	c := make(chan *protobee.Connections)
	go func() {
		for {
			message := <-c
			fmt.Println("messge received", message)
		}
	}()
	//Listen to the TCP port
	listener, err := net.Listen("tcp", "127.0.0.1:2110")
	if err != nil {
		fmt.Println("error", err)
	}
	for {
		if conn, err := listener.Accept(); err == nil {
			//If err is nil then that means that data is available for us so we take up this data and pass it to a new goroutine
			go handleProtoClient(conn, c)
		} else {
			continue
		}
	}
}
