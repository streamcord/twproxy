package config

// Service ...
type Service struct {
	Auth         string `yaml:"auth"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
}

// Config ...
type Config struct {
	Services map[string]Service `yaml:"services"`
}
