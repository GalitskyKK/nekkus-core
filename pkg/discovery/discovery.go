package discovery

import (
	"fmt"
	"strings"

	"github.com/hashicorp/mdns"
)

const ServiceTag = "_nekkus._tcp"

type ModuleAnnouncement struct {
	ID       string
	Name     string
	HTTPPort int
	GRPCPort int
	Host     string
}

func Announce(module ModuleAnnouncement) (*mdns.Server, error) {
	info := []string{
		fmt.Sprintf("id=%s", module.ID),
		fmt.Sprintf("name=%s", module.Name),
		fmt.Sprintf("grpc=%d", module.GRPCPort),
	}

	service, err := mdns.NewMDNSService(
		module.ID,
		ServiceTag,
		"",
		"",
		module.HTTPPort,
		nil,
		info,
	)
	if err != nil {
		return nil, err
	}

	return mdns.NewServer(&mdns.Config{Zone: service})
}

func Discover() ([]ModuleAnnouncement, error) {
	entriesCh := make(chan *mdns.ServiceEntry, 10)
	go func() {
		mdns.Lookup(ServiceTag, entriesCh)
		close(entriesCh)
	}()

	var modules []ModuleAnnouncement
	for entry := range entriesCh {
		module := parseEntry(entry)
		if module != nil {
			modules = append(modules, *module)
		}
	}
	return modules, nil
}

func parseEntry(entry *mdns.ServiceEntry) *ModuleAnnouncement {
	m := &ModuleAnnouncement{
		HTTPPort: entry.Port,
		Host:     entry.Host,
	}

	for _, field := range entry.InfoFields {
		parts := strings.SplitN(field, "=", 2)
		if len(parts) != 2 {
			continue
		}
		switch parts[0] {
		case "id":
			m.ID = parts[1]
		case "name":
			m.Name = parts[1]
		case "grpc":
			fmt.Sscanf(parts[1], "%d", &m.GRPCPort)
		}
	}

	if m.ID == "" {
		return nil
	}
	return m
}
