package gocord

import "testing"

func TestBasicAuth(t *testing.T) {
	var username = "stitch"
	var password = "is_awesome"
	var expected = "c3RpdGNoOmlzX2F3ZXNvbWU="

	if basicAuth(username, password) != expected {
		t.Fail()
	}
}
