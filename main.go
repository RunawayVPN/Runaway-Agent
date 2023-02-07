package main

import (
	"encoding/json"
	"os"

	"github.com/RunawayVPN/Runaway-Agent/internal/handlers"
	"github.com/RunawayVPN/Runaway-Agent/tools/hub"

	"github.com/gin-gonic/gin"
)

func main() {
	// Check if cache file exists
	hub_info, err := load_hub_info()
	if err != nil {
		panic(err)
	}
	// Open up a port to handle requests from the hub
	agent_server := gin.Default()
	// Middleware for JWT authentication
	agent_server.Use(handlers.JWTAuth(hub_info))
	// Add routes
	agent_server.POST("/add_config", handlers.AddConfig)

}

// The Hub information contains the public key and JWT authorization for the agent
// The public key is used to authenticate requests coming from the hub
// The JWT is used to authenticate requests going to the hub for additional information
func load_hub_info() (hub.HubInfo, error) {
	// Check if ~/.config/runawayvpn/cache/hub_info.json exists
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return hub.HubInfo{}, err
	}
	if _, err := os.Stat(home_dir + ".config/runawayvpn/cache/hub_info.json"); err != nil {
		// Need to register
		hub_info, err := hub.Register()
		if err != nil {
			return hub.HubInfo{}, err
		}
		return hub_info, nil
	}
	// Read hub_info from cache
	hub_info_json, err := os.ReadFile(home_dir + ".config/runawayvpn/cache/hub_info.json")
	if err != nil {
		return hub.HubInfo{}, err
	}
	hub_info := hub.HubInfo{}
	err = json.Unmarshal(hub_info_json, &hub_info)
	if err != nil {
		return hub.HubInfo{}, err
	}
	return hub_info, nil
}
