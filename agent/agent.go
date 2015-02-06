package agent

import (
	"github.com/alicebob/procspy"
	"github.com/golang/protobuf/proto"
	"github.com/infinitystrip/honeybee/protobee"
	"log"
	"net"
	"os"
	"time"
)

func sendDataToDest(data []byte, dst *string) {
	conn, err := net.Dial("tcp", *dst)
	if err != nil {
		log.Fatal("error", err)
	}
	n, err := conn.Write(data)
	if err != nil {
		log.Println("error", err)
	}
	log.Printf("Sent %d bytes\n", n)
}

func listener(c <-chan *procspy.ConnIter, dst, apiKey, label *string) {
	for {
		log.Println("Listen for new message")
		connIter := <-c

		protobeeMsg := new(protobee.Message)
		server := new(protobee.Server)
		hostname, _ := os.Hostname()
		server.Hostname = &hostname
		server.Label = label

		err := copy(connIter, &server.Connections)
		if err != nil {
			log.Println("error", err)
		}

		protobeeMsg.Type = protobee.Message_SERVER.Enum()
		protobeeMsg.Server = server
		//connections
		pb, err := proto.Marshal(protobeeMsg)
		if err != nil {
			log.Println("error", err)
		}
		sendDataToDest(pb, dst)
	}
}

func copy(connIter *procspy.ConnIter, connections *[]*protobee.Connection) error {

	conn := (*connIter).Next()

	for conn != nil {
		log.Println(". connection: ", conn)
		protoConn := new(protobee.Connection)
		protoConn.Transport = proto.String(conn.Transport)
		protoConn.LocalAddress = proto.String(conn.LocalAddress.String())
		protoConn.LocalPort = proto.Uint32(uint32(conn.LocalPort))
		protoConn.RemoteAddress = proto.String(conn.RemoteAddress.String())
		protoConn.RemotePort = proto.Uint32(uint32(conn.RemotePort))
		protoConn.Pid = proto.Uint32(uint32(conn.PID))
		protoConn.Name = proto.String(conn.Name)

		(*connections) = append((*connections), protoConn)
		conn = (*connIter).Next()
	}

	return nil
}

func startMonitor(channel chan<- *procspy.ConnIter, scanningSeconds int64) {
	for {
		cs, err := procspy.Connections(true)
		if err != nil {
			panic(err)
		}
		channel <- &cs
		time.Sleep(time.Second * time.Duration(scanningSeconds))
	}

}

func sendSystemInfo(dst *string) {
	systemInfo, _ := Collect()
	log.Println("SystemInfo", systemInfo)
	protobeeMsg := new(protobee.Message)

	protobeeMsg.Type = protobee.Message_SYSTEM_INFO.Enum()
	protobeeMsg.SystemInfo = systemInfo
	//connections
	pb, err := proto.Marshal(protobeeMsg)
	if err != nil {
		log.Println("error", err)
	}
	sendDataToDest(pb, dst)
}

func Run() {
	log.Println("Start agent")

	var c chan *procspy.ConnIter = make(chan *procspy.ConnIter)

	dst := "127.0.0.1:2110"
	apiKey := "key"
	name := "name"

	sendSystemInfo(&dst)

	go listener(c, &dst, &apiKey, &name)

	startMonitor(c, 1)
}
