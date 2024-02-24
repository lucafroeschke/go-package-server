package config

import (
	"github.com/lucafroeschke/go-package-server/logger"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"reflect"
	"sync"
)

const FileName = "config.yaml"

var (
	config        *Config
	once          sync.Once
	defaultConfig = Config{
		Domain:           "go.domain.com",
		ListeningAddress: "0.0.0.0",
		ListeningPort:    8080,
		LogRequests:      false,
		Packages: []Package{
			{
				Name:        "example",
				Repository:  "https://github.com/repository/example",
				Description: "An example package",
				Vcs:         "git",
			},
		},
	}
)

type Config struct {
	Domain           string    `yaml:"domain"`
	ListeningAddress string    `yaml:"listening_address"`
	ListeningPort    int       `yaml:"listening_port"`
	LogRequests      bool      `yaml:"log_requests"`
	Packages         []Package `yaml:"packages"`
}

type Package struct {
	Name        string `yaml:"name"`
	Repository  string `yaml:"repository"`
	Description string `yaml:"description"`
	Vcs         string `yaml:"vcs"`
}

func GetConfig() *Config {
	once.Do(func() {
		data, err := os.ReadFile(FileName)
		if err != nil {
			if os.IsNotExist(err) {
				config = &defaultConfig

				logger.WriteLog(logger.INFO, "Creating new config file")
				err := SaveConfig()
				if err != nil {
					log.Fatalf("Failed to create config file: %v", err)
				}
			} else {
				log.Fatalf("Failed to read config file: %v", err)
			}
		} else {
			err = yaml.Unmarshal(data, &config)

			if err != nil {
				log.Fatalf("Failed to unmarshal config data: %v", err)
			}

			if checkMissingConfigFields(config) {
				logger.WriteLog(logger.INFO, "Added missing fields to config")
				err := SaveConfig()
				if err != nil {
					log.Fatalf("Failed to save config file: %v", err)
				}
			}

			logger.WriteLog(logger.INFO, "Loaded config file")
		}
	})

	return config
}

func GetPackage(name string) (*Package, bool) {
	cfg := GetConfig()
	for _, p := range cfg.Packages {
		if p.Name == name {
			return &p, true
		}
	}

	return nil, false
}

func SaveConfig() error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile(FileName, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write config file: %v", err)
		return err
	}

	return nil
}

func checkMissingConfigFields(config *Config) bool {
	updated := false
	configType := reflect.TypeOf(*config)
	configValue := reflect.ValueOf(config).Elem()
	defaultConfigValue := reflect.ValueOf(defaultConfig)
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		fieldValue := configValue.Field(i)
		if reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(field.Type).Interface()) {
			defaultFieldValue := defaultConfigValue.FieldByName(field.Name)
			fieldValue.Set(defaultFieldValue)
			updated = true
		}
	}
	return updated
}
