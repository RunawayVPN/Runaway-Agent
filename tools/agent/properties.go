package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/RunawayVPN/security"
	"github.com/RunawayVPN/types"
)

type network_details struct {
	PublicIP string `json:"public_ip"`
	Country  string `json:"country"`
	ISP      string `json:"isp"`
}

func fetch_network_details() (network_details, error) {
	// HTTP GET request to http://ip-api.com/json/
	response, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return network_details{}, err
	}
	defer response.Body.Close()
	// Read response body
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	// Decode response body as JSON
	var response_body map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &response_body)
	if err != nil {
		return network_details{}, err
	}
	// Construct interface
	return network_details{
		PublicIP: response_body["query"].(string),
		Country:  response_body["country"].(string),
		ISP:      response_body["isp"].(string),
	}, nil
}

func Construct_agent() (types.Agent, error) {
	// Fetch network details
	network_details, err := fetch_network_details()
	if err != nil {
		return types.Agent{}, err
	}
	// Fetch public key
	pubkey := security.EncodeBS(security.Public_key)
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		// Generate random name with 8 characters
		hostname = create_random_name(8)
	}
	// Construct agent
	return types.Agent{
		PublicIP:  network_details.PublicIP,
		Country:   network_details.Country,
		ISP:       network_details.ISP,
		PublicKey: pubkey,
		Name:      hostname,
	}, nil
}

func create_random_name(length int) string {
	// Create random name
	var name string
	for i := 0; i < length; i++ {
		name += fmt.Sprintf("%c", rand.Intn(26)+65)
	}
	return name
}
