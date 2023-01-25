package hub

import (
	"bytes"
	"encoding/json"
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

func Register() error {
	// Get Agent
	agent, err := agent.Construct_agent()
	if err != nil {
		return err
	}
	// Convert agent to JSON string
	agent_json, err := json.Marshal(agent)
	if err != nil {
		return err
	}
	agent_jwt, err := security.CreateToken(string(agent_json))
	if err != nil {
		return err
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
		return err
	}
	// Make request to endpoint
	resp, err := http.Post(endpoint, "application/json", bytes.NewBuffer([]byte(registration_request)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Map response to RegistrationResponse struct
	var registration_response types.RegistrationResponse
	err = json.NewDecoder(resp.Body).Decode(&registration_response)
	if err != nil {
		return err
	}
	// Verify JWT
	_, err = security.VerifyToken(registration_response.JwtToken, registration_response.PublicKey)
	if err != nil {
		return err
	}
	// Save public key to environment variable
	err = os.Setenv("HUB_PUBLIC_KEY", registration_response.PublicKey)
	if err != nil {
		return err
	}
	err = os.Setenv("AGENT_JWT", registration_response.JwtToken)
	if err != nil {
		return err
	}
	// Save
	return nil
}
