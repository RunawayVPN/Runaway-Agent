package hub

import "testing"

func TestRegistration(t *testing.T) {
	err := Register()
	if err != nil {
		t.Error(err)
	}
}
