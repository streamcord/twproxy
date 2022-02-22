package config

// Service ...
type Service struct {
	Auth         string `yaml:"auth"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

// DogStatsD ...
type DogStatsD struct {
	Address   string `yaml:"address"`
	Namespace string `yaml:"namespace"`
}

// Config ...
type Config struct {
	DogStatsD DogStatsD          `yaml:"dogStatsD"`
	Services  map[string]Service `yaml:"services"`
}
