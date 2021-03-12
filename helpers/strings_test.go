package helpers

import "testing"

func TestSha256String(t *testing.T) {
	dataItems := []struct {
		Data   string
		Result string
	}{
		{"abdulsamet", "fc1b1cb8b3494ad3a2142300143b6c9dd583aa3a98d38bdf96448ed63957b399"},
		{"ileri", "40a3e110e0d9dd7b2667b50c90f4f585e7630437748e40e2ae9fbf9ffbf4923d"},
	}
	for _, item := range dataItems {
		if hash := Sha256String(item.Data); hash != item.Result {
			t.Errorf("Expected: %v, Arriving: %v", item.Result, hash)
		}
	}
}
