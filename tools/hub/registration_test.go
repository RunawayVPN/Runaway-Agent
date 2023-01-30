package hub

import "testing"

func TestRegistration(t *testing.T) {
	hub_info, err := Register()
	if err != nil {
		t.Error(err)
	}
	t.Log(hub_info)
}
