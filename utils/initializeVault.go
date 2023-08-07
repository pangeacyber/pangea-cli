package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

type Config struct {
	Title    string             `mapstructure:"title"`
	Services map[string]Service `mapstructure:"services"`
}

type Service struct {
	PANGEA_DOMAIN string `mapstructure:"PANGEA_DOMAIN"`
	PANGEA_TOKEN  string `mapstructure:"PANGEA_TOKEN"`
}

func CreateVaultAPIClient() *resty.Client {
	var token string
	var err error

	defaultToken := os.Getenv("PANGEA_TOKEN")
	// Ignore reading token from file if token given through the PANGEA_TOKEN env variable
	if defaultToken == "" {
		token, err = readTokenFromConfig()
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		token = defaultToken
	}

	client := resty.New()
	client.SetAuthToken(token)
	client.SetHeader("Content-Type", "application/json")

	return client
}

// ReadTokenFromConfig reads the token from the ~/.pangea/config file.
// If the file or folder doesn't exist, it returns an error.
func readTokenFromConfig() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Create the .pangea folder if it doesn't exist
	pangeaDir := filepath.Join(homeDir, ".pangea")
	if err != nil {
		return "", err
	}

	// Create or open the config file
	configPath := filepath.Join(pangeaDir, "config.toml")

	configViper := viper.New()
	configViper.SetConfigFile(configPath)
	// Check if the config file exists
	if _, err := os.Stat(configPath); err != nil {
		fmt.Println(err)
		return "", fmt.Errorf("Pangea Token doesn't exist. Run `pangea login` to setup your CLI.")
	}

	// Read the configuration
	if err := configViper.ReadInConfig(); err != nil {
		return "", fmt.Errorf("Error fetching your pangea token: %s", err)
	}

	// Get the token from the configuration
	var config Config
	if err := configViper.Unmarshal(&config); err != nil {
		return "", fmt.Errorf("Error unmarshaling config: %s", err)
	}

	token := config.Services["default"].PANGEA_TOKEN

	if token == "" {
		return "", fmt.Errorf("Pangea Token doesn't exist. Run `pangea login` to setup your CLI.")
	}

	return token, nil
}

// writeTokenToFile writes the Pangea token to the specified file
func WriteTokenToFile(token string, domain string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// Create the .pangea folder if it doesn't exist
	pangeaDir := filepath.Join(homeDir, ".pangea")
	err = os.MkdirAll(pangeaDir, 0700)
	if err != nil {
		return err
	}

	// Create or open the config file
	configPath := filepath.Join(pangeaDir, "config.toml")
	configViper := viper.New()
	configViper.SetConfigFile(configPath)

	newConfig := Config{
		Title: "Pangea",
		Services: map[string]Service{
			"default": {
				PANGEA_TOKEN:  token,
				PANGEA_DOMAIN: domain,
			},
		},
	}

	// Write the new config to the TOML file
	configViper.Set("title", newConfig.Title)
	configViper.Set("services", newConfig.Services)
	if err := configViper.WriteConfig(); err != nil {
		return err
	}

	fmt.Printf("Token successfully written to %s\n", configPath)
	return nil
}
