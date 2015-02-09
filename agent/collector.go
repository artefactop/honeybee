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
		log.Println("collector:", c)
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
		case "filesystem":
			sliceFs := c.([]interface{})

			for _, v := range sliceFs {
				fs := v.(map[string]string)
				protoFs := new(protobee.FileSystem)

				protoFs.Name = proto.String(fs["name"])
				protoFs.KbSize = mapAtoi(fs["kb_size"])
				protoFs.MountedOn = proto.String(fs["mounted_on"])

				systemInfo.FileSystems = append(systemInfo.FileSystems, protoFs)
			}
		case "memory":
			mem := c.(map[string]string)
			log.Println("mem", mem)
			systemInfo.Memory = new(protobee.Memory)

			systemInfo.Memory.SwapTotal = proto.String(mem["swap_total"])
			systemInfo.Memory.Total = mapAtoi(mem["total"]) //use uint64
		case "network":
			net := c.(map[string]interface{}) //TODO get all network interfaces
			log.Println("net", net)

			protoNet := new(protobee.Network)
			protoNet.MacAddress = proto.String(net["macaddress"].(string))
			protoNet.IpAddress = proto.String(net["ipaddress"].(string))
			protoNet.IpAddressV6 = proto.String(net["ipaddressv6"].(string))

			systemInfo.Networks = append(systemInfo.Networks, protoNet)

		case "platform":
			platform := c.(map[string]interface{})
			log.Println("platform", platform)
			systemInfo.Platform = new(protobee.Platform)

			systemInfo.Platform.Hostname = proto.String(platform["hostname"].(string))
			systemInfo.Platform.KernelName = proto.String(platform["kernel_name"].(string))
			systemInfo.Platform.KernelRelease = proto.String(platform["kernel_release"].(string))
			systemInfo.Platform.Machine = proto.String(platform["machine"].(string))
			systemInfo.Platform.Processor = proto.String(platform["processor"].(string))
			systemInfo.Platform.Os = proto.String(platform["os"].(string))

		default:
			log.Println("collector", c)
		}
	}

	return
}
