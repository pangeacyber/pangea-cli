package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-resty/resty/v2"
)

func CreateVaultAPIClient() *resty.Client {
	token, err := readTokenFromConfig()
	if err != nil {
		log.Fatalln(err)
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
	configPath := filepath.Join(pangeaDir, "config")

	tokenBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("Error fetching your pangea token: %s", err)
	}

	token := strings.TrimSpace(string(tokenBytes))
	if token == "" {
		return "", fmt.Errorf("Pangea Token doesn't exist. Run pangea login to setup your CLI.")
	}

	return token, nil
}

// writeTokenToFile writes the Pangea token to the specified file
func WriteTokenToFile(token string) error {
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
	configPath := filepath.Join(pangeaDir, "config")
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(token)
	if err != nil {
		return err
	}

	fmt.Printf("Token successfully written to %s\n", configPath)
	return nil
}
