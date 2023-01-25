package agent

import "testing"

func TestAgentConstruction(t *testing.T) {
	agent, err := Construct_agent()
	if err != nil {
		t.Error(err)
	}
	t.Log(agent)
}
