package domain

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Kraken Kraken `yaml:"kraken"`
	Smtp   Smtp   `yaml:"smtp"`

	Notify    string    `yaml:"notify"`
	Frequency string    `yaml:"frequency"`
	Currency  string    `yaml:"currency"`
	Pairs     []DCAPair `yaml:"pairs"`
}

type Kraken struct {
	Key    string `yaml:"key"`
	Secret string `yaml:"secret"`
}

type Smtp struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
}

type DCAPair struct {
	Pair   string  `json:"pair"`
	Amount float64 `json:"amount"`
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while trying to open the configuration file : %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Print("An error occurred while closing the configuration file")
		}
	}(file)

	yamlString, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("cannot read the configuration file : %w", err)
	}
	yamlString = []byte(os.ExpandEnv(string(yamlString)))

	var config Config
	err = yaml.Unmarshal(yamlString, &config)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while unmarshalling the configuration Yaml : %w", err)
	}

	if config.Kraken.Key == "" {
		return nil, errors.New("the kraken key is not specified")
	}

	if config.Kraken.Secret == "" {
		return nil, errors.New("the kraken secret is not specified")
	}

	return &config, nil
}
