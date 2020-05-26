package user

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "123456"
	if hash, err := HashPassword(password); err != nil {
		t.Error(err)
	} else {
		fmt.Println("Hash:", hash)
		match := CheckPasswordHash(password, hash)
		fmt.Println("Match:   ", match)
	}
}
