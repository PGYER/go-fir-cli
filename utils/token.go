package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func LoadLocalToken() string {
	config, err := LoadLocalConfig()
	if err != nil {
		return ""
	}
	return config["token"]
}

func LoadLocalConfig() (map[string]string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	tokenPath := filepath.Join(home, ".fir-cli")
	raw, err := os.ReadFile(tokenPath)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(string(raw), "---\n") {
		raw = raw[4:]
	}
	data := make(map[string]string)
	for _, item := range strings.Split(string(raw), "\n") {
		cmp := strings.SplitN(item, ": ", 2)
		if len(cmp) == 2 && strings.HasPrefix(cmp[0], ":") {
			cmp[0] = cmp[0][1:len(cmp[0])]
			data[cmp[0]] = cmp[1]
		}
	}
	return data, nil
}

func SaveToLocal(email string, token string) error {
	raw := fmt.Sprintf("---\n:email: %s\n:token: %s\n", email, token)
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	tokenPath := filepath.Join(home, ".fir-cli")
	return os.WriteFile(tokenPath, []byte(raw), 0644)
}
