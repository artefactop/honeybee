package main

import (
	"fmt"
	"github.com/infinitystrip/honeybee/server"
)

//TODO save data in influxdb

func main() {
	nodes := make(map[string]string)
	edges := make(map[string][]string)
	fmt.Println("Starting honeybee server")
	go server.ListenPB(&nodes, &edges)
	server.Run(&nodes, &edges)
}
