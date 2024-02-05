package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type CacheData struct {
	Paths map[string]ProjectData `json:"paths"`
}

type ProjectData struct {
	Remote string `json:"remote"`
}

func CheckPathExists() (bool, CacheData, string) {
	cachePath := GetCachePath()
	config := LoadCacheData(cachePath)

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Error reading current directory path")
	}

	currentDir = strings.ToLower(currentDir)

	if _, isPathExists := config.Paths[currentDir]; isPathExists {
		return true, config, currentDir
	} else {
		return false, config, currentDir
	}
}

func GetCachePath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".pangea", "cache_paths.json")
}

func LoadCacheData(cachePath string) CacheData {
	internalViper := viper.New()
	internalViper.SetConfigFile(cachePath)
	internalViper.SetConfigType("json")
	internalViper.SetDefault("paths", make(map[string]ProjectData))

	internalViper.ReadInConfig() //nolint:errcheck

	var config CacheData
	err := internalViper.Unmarshal(&config)
	if err != nil {
		fmt.Printf("Error parsing cache file: %v\n", err)
	}

	return config
}
