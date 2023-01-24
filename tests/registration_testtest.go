package main

import (
	"bytes"
	"encoding/json"
	"net/http"

	sec "github.com/RunawayVPN/Runaway-Hub/tools/security"
)

func main() {
	// Create JWT payload
	type RegistrationPayload struct {
		PublicIP  string `json:"public_ip"`
		SecretKey string `json:"secret_key"`
		Name      string `json:"name"`
	}
	var payload RegistrationPayload = RegistrationPayload{
		PublicIP:  "1.1.1.4",
		SecretKey: "secret",
		Name:      "test",
	}
	// Make payload JSON string
	json_payload, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	jwt, err := sec.CreateToken(string(json_payload))
	if err != nil {
		panic(err)
	}
	// Make HTTP POST request
	body := map[string]string{
		"jwt":        jwt,
		"public_key": sec.EncodeBS(sec.Public_key),
	}
	// Encode body for HTTP POST request (JSON)
	json_body, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	response, err := http.Post("http://localhost:8080/agent/registration", "application/json", bytes.NewBuffer(json_body))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	// Print response
	println(response.Status)
	// Print response body
	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	println(buf.String())
	// Read response body as JSON
	var response_body map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &response_body)
	if err != nil {
		panic(err)
	}
	// Verify JWT
	jwt_response, err := sec.VerifyToken(response_body["jwt"].(string), response_body["public_key"].(string))
	if err != nil {
		panic(err)
	}
	println(jwt_response)
}
