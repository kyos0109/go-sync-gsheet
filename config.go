package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"gopkg.in/yaml.v2"
)

type yamlConfig struct {
	Setting setting   `yaml:"Setting"`
	AWS     awsConfig `yaml:"AWS"`
}

type setting struct {
	SyncTimeInterval time.Duration `yaml:"SyncTimeInterval"`
	SpreadsheetId    string        `yaml:"SpreadsheetId"`
}

type awsConfig struct {
	Auth []awsAuth `yaml:"Auth"`
}

type awsAuth struct {
	Account   string `yaml:"Account"`
	AccessKey string `yaml:"AccessKey"`
	SecretKey string `yaml:"SecretKey"`
	Region    string `yaml:"Region"`
	Project   string `yaml:"Project,omitempty"`
}

func GetConfig(configPath string) *yamlConfig {
	var yc yamlConfig
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Println(err)
	}

	err = yaml.Unmarshal(yamlFile, &yc)
	if err != nil {
		fmt.Println("Unmarshal:", err)
	}

	return &yc
}
