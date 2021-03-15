package helpers

import "testing"

func TestSha256String(t *testing.T) {
	testCases := []struct {
		Data   string
		Result string
	}{
		{"abdulsamet", "fc1b1cb8b3494ad3a2142300143b6c9dd583aa3a98d38bdf96448ed63957b399"},
		{"ileri", "40a3e110e0d9dd7b2667b50c90f4f585e7630437748e40e2ae9fbf9ffbf4923d"},
	}
	for _, test := range testCases {
		if hash := Sha256String(test.Data); hash != test.Result {
			t.Errorf("Expected: %v, Want: %v", test.Result, hash)
		}
	}
}

func TestStripBearer(t *testing.T) {
	testCases := []struct {
		Token  string
		Result string
	}{
		{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcl9uYW1lIjoiYWJkdWxzYW1ldCIsImV4cCI6MTY0NzI3Njg0N30.ja-BZSyGLi5Kkj4oBxiGm-s0PK-aGKj2aJvXMabFjk0", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcl9uYW1lIjoiYWJkdWxzYW1ldCIsImV4cCI6MTY0NzI3Njg0N30.ja-BZSyGLi5Kkj4oBxiGm-s0PK-aGKj2aJvXMabFjk0"},
		{"beaRer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcl9uYW1lIjoiYWJkdWxzYW1ldCIsImV4cCI6MTY0NzI3Njg0N30.ja-BZSyGLi5Kkj4oBxiGm-s0PK-aGKj2aJvXMabFjk0", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwidXNlcl9uYW1lIjoiYWJkdWxzYW1ldCIsImV4cCI6MTY0NzI3Njg0N30.ja-BZSyGLi5Kkj4oBxiGm-s0PK-aGKj2aJvXMabFjk0"},
		{"testtt", "testtt"},
	}
	for _, test := range testCases {
		if token, _ := StripBearer(test.Token); token != test.Result {
			t.Errorf("Expected: %v, Want: %v", test.Result, token)
		}
	}
}
