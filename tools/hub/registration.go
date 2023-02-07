package hub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/RunawayVPN/Runaway-Agent/tools/agent"
	"github.com/RunawayVPN/security"
	"github.com/RunawayVPN/types"
)

// Private types
type HubInfo struct {
	PublicKey string
	AuthToken string
}

// Global variables
var path string
var endpoint string

func init() {
	// Get endpoint environment variable
	protocol := os.Getenv("PROTOCOL")
	endpoint = os.Getenv("ENDPOINT")
	port := os.Getenv("PORT")
	if endpoint == "" {
		endpoint = "localhost"
	}
	path = protocol + endpoint + port
}

func Register() (HubInfo, error) {
	// Get Agent
	agent, err := agent.Construct_agent()
	if err != nil {
		return HubInfo{}, err
	}
	// Get secret key from environment variable
	secret_key := os.Getenv("SECRET_KEY")
	if secret_key == "" {
		secret_key = "secret"
	}
	// Generate AuthToken based on endpoint
	auth_token_json, err := json.Marshal(types.AuthToken{
		Endpoint: endpoint,
		Roles:    []string{"hub"},
	})
	if err != nil {
		return HubInfo{}, err
	}
	auth_token, err := security.CreateToken(string(auth_token_json))
	if err != nil {
		return HubInfo{}, err
	}
	request, err := json.Marshal(types.RegistrationRequest{
		PublicKey: agent.PublicKey,
		SecretKey: secret_key,
		Agent:     agent,
		AuthToken: auth_token,
	})
	if err != nil {
		return HubInfo{}, err
	}
	// Make request to endpoint
	resp, err := http.Post(path+"/agent/registration", "application/json", bytes.NewBuffer([]byte(request)))
	if err != nil {
		return HubInfo{}, err
	}
	defer resp.Body.Close()
	// Map response to RegistrationResponse struct
	var response types.RegistrationResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return HubInfo{}, err
	}
	if resp.StatusCode != 200 {
		return HubInfo{}, fmt.Errorf("registration failed: %s", response.Error)
	}
	// Verify JWT
	_, err = security.VerifyToken(response.AuthToken, response.PublicKey)
	if err != nil {
		return HubInfo{}, err
	}
	// Return HubInfo
	return HubInfo{
		PublicKey: response.PublicKey,
		AuthToken: response.AuthToken,
	}, nil
}
