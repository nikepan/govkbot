package govkbot

import (
	"testing"
)

func TestUser_FullName(t *testing.T) {
	var u User
	if u.FullName() != "" {
		t.Error("User name must be blank")
	}
	u = User{FirstName: "First", LastName: "Last"}
	if u.FullName() != "First Last" {
		t.Error("Wrong full user name")
	}
}
