package govkbot

import (
	"testing"
)

func TestUser_FullName(t *testing.T) {
	u := User{FirstName: "First", LastName: "Last"}
	if u.FullName() != "First Last" {
		t.Error("Wrong full user name")
	}
}