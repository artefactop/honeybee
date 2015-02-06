package agent

import (
	"github.com/DataDog/gohai/cpu"
	"github.com/DataDog/gohai/filesystem"
	"github.com/DataDog/gohai/memory"
	"github.com/DataDog/gohai/network"
	"github.com/DataDog/gohai/platform"
	"github.com/golang/protobuf/proto"
	"github.com/infinitystrip/honeybee/protobee"
	"log"
	"strconv"
)

type Collector interface {
	Name() string
	Collect() (interface{}, error)
}

var collectors = []Collector{
	&cpu.Cpu{},
	&filesystem.FileSystem{},
	&memory.Memory{},
	&network.Network{},
	&platform.Platform{},
}

func mapAtoi(org string) *uint32 {
	i, err := strconv.Atoi(org)
	if err != nil {
		log.Printf("[%s] %s", org, err)
	}
	return proto.Uint32(uint32(i))
}

func Collect() (systemInfo *protobee.SystemInfo, err error) {
	systemInfo = new(protobee.SystemInfo)

	for _, collector := range collectors {
		c, err := collector.Collect()
		if err != nil {
			log.Printf("[%s] %s", collector.Name(), err)
			continue
		}
		switch collector.Name() {
		case "cpu":
			cpu := c.(map[string]string)
			systemInfo.Cpu = new(protobee.Cpu)

			systemInfo.Cpu.CpuCores = mapAtoi(cpu["cpu_cores"])
			systemInfo.Cpu.Family = mapAtoi(cpu["family"])
			systemInfo.Cpu.Mhz = mapAtoi(cpu["mhz"])
			systemInfo.Cpu.Model = mapAtoi(cpu["model"])
			systemInfo.Cpu.Stepping = mapAtoi(cpu["stepping"])
			systemInfo.Cpu.ModelName = proto.String(cpu["model_name"])
			systemInfo.Cpu.VendorId = proto.String(cpu["vendor_id"])
		default:
			log.Println("collector", c)
		}
	}

	return
}
