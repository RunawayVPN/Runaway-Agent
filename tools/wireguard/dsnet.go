package wireguard

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	dsnet "github.com/RunawayVPN/dsnet/cmd/cli"
)

func init() {
	// Check if dsnet configuration exists
	_, err := dsnet.MustLoadConfigFile()
	if err != nil {
		dsnet.Init()
		// Necessary edits:
		// "PostUp": "iptables -A FORWARD -i %i -j ACCEPT; iptables -A FORWARD -o %i -j ACCEPT; iptables -t nat -A POSTROUTING -o <main network interface> -j MASQUERADE",
		// "PostDown": "iptables -D FORWARD -i %i -j ACCEPT; iptables -D FORWARD -o %i -j ACCEPT; iptables -t nat -D POSTROUTING -o <main network interface> -j MASQUERADE",
		// 	"Networks": [
		//     "0.0.0.0/0",
		//     "::/0"
		// ],

		// Read /etc/dsnetconfig.json
		config, err := os.ReadFile("/etc/dsnetconfig.json")
		if err != nil {
			panic(err)
		}
		// Parse JSON
		var data map[string]interface{}
		err = json.Unmarshal(config, &data)
		if err != nil {
			panic(err)
		}
		// Get main network interface
		main_interface, err := get_main_network_interface()
		if err != nil {
			panic(err)
		}
		// Edit PostUp and PostDown
		data["PostUp"] = fmt.Sprintf("iptables -A FORWARD -i %%i -j ACCEPT; iptables -A FORWARD -o %%i -j ACCEPT; iptables -t nat -A POSTROUTING -o %s -j MASQUERADE", main_interface)
		data["PostDown"] = fmt.Sprintf("iptables -D FORWARD -i %%i -j ACCEPT; iptables -D FORWARD -o %%i -j ACCEPT; iptables -t nat -D POSTROUTING -o %s -j MASQUERADE", main_interface)
		// Edit Networks
		data["Networks"] = []string{"0.0.0.0/0", "::/0"}
		// Write to /etc/dsnetconfig.json
		config, err = json.MarshalIndent(data, "", "  ")
		if err != nil {
			panic(err)
		}
		os.WriteFile("/etc/dsnetconfig.json", config, 0644)
	}
}

// Get the main network interface
func get_main_network_interface() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, i := range interfaces {
		if i.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if i.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return i.Name, nil
		}
	}
	return "", fmt.Errorf("network interface not found")
}
