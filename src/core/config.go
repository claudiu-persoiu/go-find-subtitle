package core

import (
	"encoding/json"
	"fmt"
	"github.com/claudiu-persoiu/go-find-subtitle/src/provider/opensubtitles"
	"os"
)

type Config struct {
	Opensubititles opensubtitles.Config `json:"opensubititles"`
}

func GetConfig() (*Config, error) {
	raw, err := os.ReadFile("/etc/go-find-subtitle.json")
	if err != nil {
		return nil, fmt.Errorf("error occured while reading config: %s", err)
	}
	var config Config
	err = json.Unmarshal(raw, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmashalling config: %s", err)
	}

	return &config, nil
}
