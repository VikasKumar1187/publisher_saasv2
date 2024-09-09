package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vikaskumar1187/publisher_saasv2/services/publisher/foundation/resty"
)

type Environment string

const (
	Dev   Environment = "dev"
	Stage Environment = "stage"
	Prod  Environment = "prod"
)

type CustomerConfig struct {
	ClientID     string `json:"clientID"`
	ClientSecret string `json:"clientSecret"`
	StoryURL     string `json:"storyUrl"`
}

type Config struct {
	TokenURLs map[Environment]string    `json:"TokenURLs"`
	Customers map[string]CustomerConfig `json:"Customers"`
}

type TokenRequest struct {
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <customer> <environment> <action>", os.Args[0])
	}

	customer := os.Args[1]
	env := Environment(os.Args[2])
	action := os.Args[3]

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	switch action {
	case "tokengen":
		if err := generateToken(config, customer, env); err != nil {
			log.Fatalf("Error generating token: %v", err)
		}
	default:
		log.Fatalf("Unknown action: %s", action)
	}
}

func loadConfig() (*Config, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("error getting current working directory: %w", err)
	}
	log.Printf("Current working directory: %s", cwd)

	configPath := filepath.Join("app", "tooling", "admin", "auth", "config.json")
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %w", err)
	}
	log.Printf("Attempting to load config from: %s", absConfigPath)

	file, err := os.Open(configPath)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}

	// Debug: Print loaded configuration
	log.Printf("Loaded configuration: %+v", config)

	return &config, nil
}

func generateToken(config *Config, customerName string, env Environment) error {
	// Debug: Print available customers
	log.Printf("Available customers: %v", getCustomerNames(config))

	customerConfig, ok := findCustomer(config, customerName)
	if !ok {
		return fmt.Errorf("customer %s not found", customerName)
	}

	tokenURL, ok := config.TokenURLs[env]
	if !ok {
		return fmt.Errorf("token URL for environment %s not found", env)
	}

	log.Printf("Generating token for customer %s in environment %s", customerName, env)
	log.Printf("Using token URL: %s", tokenURL)

	client := resty.New(tokenURL, 10*time.Second)

	tokenReq := TokenRequest{
		GrantType:    "client_credentials",
		ClientID:     customerConfig.ClientID,
		ClientSecret: customerConfig.ClientSecret,
	}

	payload, err := resty.NewJSONPayload(tokenReq)
	if err != nil {
		return fmt.Errorf("error creating payload: %w", err)
	}

	var tokenResp TokenResponse
	err = client.Post(context.Background(), "", nil, payload, &tokenResp)
	if err != nil {
		return fmt.Errorf("failed to generate token: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return fmt.Errorf("received empty access token")
	}

	err = os.Setenv("TOKEN", tokenResp.AccessToken)
	if err != nil {
		return fmt.Errorf("error setting TOKEN environment variable: %w", err)
	}

	log.Printf("Token generated successfully and stored in TOKEN environment variable")
	return nil
}

// Helper function to find a customer case-insensitively
func findCustomer(config *Config, customerName string) (CustomerConfig, bool) {
	for name, cfg := range config.Customers {
		if strings.EqualFold(name, customerName) {
			return cfg, true
		}
	}
	return CustomerConfig{}, false
}

// Helper function to get a list of customer names
func getCustomerNames(config *Config) []string {
	var names []string
	for name := range config.Customers {
		names = append(names, name)
	}
	return names
}
