package utils

import (
	"fmt"
	"testing"
)

func TestLoadLocalToken(t *testing.T) {
	for i := 0; i < 2; i++ {
		data, err := LoadLocalConfig()
		if err != nil {
			t.Error(err)
		}
		fmt.Println(data)
		err = SaveToLocal(data["email"], data["token"])
		if err != nil {
			t.Error(err)
		}
	}
}
