package config

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

// FromFile loads a rules config file
func FromFile(filename string) (*RulesConfig, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, fmt.Errorf("unable to load rules config file: %v", err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("unable to load rules config file: %v", err)
	}

	var cfg RulesConfig
	if err := yaml.UnmarshalStrict(content, &cfg); err != nil {
		return nil, fmt.Errorf("unable to parse rules config: %v", err)
	}
	return &cfg, nil
}
