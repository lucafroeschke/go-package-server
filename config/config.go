package config

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/lucafroeschke/go-package-server/args"
	"github.com/lucafroeschke/go-package-server/logger"
	"gopkg.in/yaml.v2"
)

var FileName = "config.yml"

func init() {
	configPath := args.GetString("config-path")
	if configPath != "" {
		FileName = configPath
	}

	logger.WriteLog(logger.INFO, "Using config file: "+FileName)
}

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
	Domain           string     `yaml:"domain"`
	ListeningAddress string     `yaml:"listening_address"`
	ListeningPort    int        `yaml:"listening_port"`
	LogRequests      bool       `yaml:"log_requests"`
	DisableGoDoc     bool       `yaml:"disable_go_doc"`
	SiteConfig       SiteConfig `yaml:"site_config"`
	Packages         []Package  `yaml:"packages"`
}

type Package struct {
	Name         string `yaml:"name"`
	Repository   string `yaml:"repository"`
	Description  string `yaml:"description"`
	Vcs          string `yaml:"vcs"`
	DisableGoDoc bool   `yaml:"disable_go_doc"`
}

type SiteConfig struct {
	Title       string       `yaml:"title"`
	Description string       `yaml:"description"`
	FooterLinks []FooterLink `yaml:"footer_links"`
}

type FooterLink struct {
	Title string `yaml:"title"`
	Url   string `yaml:"url"`
}

func GetConfig() *Config {
	once.Do(func() {
		dir := filepath.Dir(FileName)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		}

		data, err := os.ReadFile(FileName)
		if err != nil {
			if os.IsNotExist(err) {
				err := CreateConfig()
				if err != nil {
					return
				}
			} else {
				log.Fatalf("Failed to read config file: %v", err)
			}
		} else {
			err = yaml.Unmarshal(data, &config)

			if err != nil {
				log.Fatalf("Failed to unmarshal config data: %v", err)
			}

			logger.WriteLog(logger.INFO, "Loaded config file")
		}
	})

	return config
}

// SetConfig sets the configuration manually. This is primarily used for testing.
func SetConfig(c *Config) {
	config = c
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

func CreateConfig() error {
	logger.WriteLog(logger.INFO, "Creating config file")
	config = &defaultConfig
	return SaveConfig()
}
