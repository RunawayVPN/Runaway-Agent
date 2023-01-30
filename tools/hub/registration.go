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

// Global variables
var endpoint string

func init() {
	// Get endpoint environment variable
	endpoint = os.Getenv("ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:8080"
	}
}

func Register() (types.HubInfo, error) {
	// Get Agent
	agent, err := agent.Construct_agent()
	if err != nil {
		return types.HubInfo{}, err
	}
	// Convert agent to JSON string
	agent_json, err := json.Marshal(agent)
	if err != nil {
		return types.HubInfo{}, err
	}
	agent_jwt, err := security.CreateToken(string(agent_json))
	if err != nil {
		return types.HubInfo{}, err
	}
	// Get secret key from environment variable
	secret_key := os.Getenv("SECRET_KEY")
	if secret_key == "" {
		secret_key = "secret"
	}
	registration_request, err := json.Marshal(types.RegistrationRequest{
		PublicKey: agent.PublicKey,
		SecretKey: secret_key,
		JwtToken:  agent_jwt,
	})
	if err != nil {
		return types.HubInfo{}, err
	}
	// Make request to endpoint
	resp, err := http.Post(endpoint+"/agent/registration", "application/json", bytes.NewBuffer([]byte(registration_request)))
	if err != nil {
		return types.HubInfo{}, err
	}
	defer resp.Body.Close()
	// Map response to RegistrationResponse struct
	var registration_response types.RegistrationResponse
	err = json.NewDecoder(resp.Body).Decode(&registration_response)
	if err != nil {
		return types.HubInfo{}, err
	}
	if resp.StatusCode != 200 {
		return types.HubInfo{}, fmt.Errorf("registration failed: %s", registration_response.Error)
	}
	// Verify JWT
	_, err = security.VerifyToken(registration_response.JwtToken, registration_response.PublicKey)
	if err != nil {
		println(registration_response.JwtToken)
		return types.HubInfo{}, err
	}
	// Save hub info
	hub_info := types.HubInfo{
		PublicKey: registration_response.PublicKey,
		AgentJwt:  registration_response.JwtToken,
	}
	// Save
	return hub_info, nil
}
