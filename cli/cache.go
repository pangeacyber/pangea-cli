package cli

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type CacheData struct {
	Paths            map[string]WorkspaceData `json:"paths"`
	VersionAvailable map[string]*Release      `json:"version_available"`
}

// Key: path in lower case. Value: WorkspaceData
type Paths map[string]WorkspaceData

type WorkspaceData struct {
	Remote string `json:"remote"`
}

// Key: DaySinceEpoch. Value: Version available
type VersionAvailable map[string]Release

func getCachePath() (string, error) {
	cacheDir, err := GetCacheFolder()
	if err != nil {
		return "", err
	}
	filepath := filepath.Join(cacheDir, "cache.json")
	_ = initFile(filepath)
	return filepath, nil
}

func LoadCacheData() (*CacheData, error) {
	cachePath, err := getCachePath()
	if err != nil {
		return nil, err
	}

	return loadCacheData(cachePath)
}

func CacheGetPaths() (Paths, error) {
	config, err := LoadCacheData()
	if err != nil {
		return nil, err
	}

	return config.Paths, nil
}

func CacheSetPaths(paths Paths) error {
	cachePath, err := getCachePath()
	if err != nil {
		return err
	}

	cd, err := loadCacheData(cachePath)
	if err != nil {
		return err
	}

	cd.Paths = paths
	return cacheSave(cd)
}

func CacheGetRelease() (*Release, error) {
	config, err := LoadCacheData()
	if err != nil {
		return nil, err
	}

	r, ok := config.VersionAvailable[GetDaySinceEpoch()]
	if !ok {
		return nil, nil
	}

	return r, nil
}

func CacheSetRelease(r *Release) error {
	cachePath, err := getCachePath()
	if err != nil {
		return err
	}

	cd, err := loadCacheData(cachePath)
	if err != nil {
		return err
	}

	cd.VersionAvailable = map[string]*Release{
		GetDaySinceEpoch(): r,
	}
	return cacheSave(cd)
}

func loadCacheData(cachePath string) (*CacheData, error) {
	jsonFile, err := os.Open(cachePath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var cd CacheData
	if err = json.Unmarshal(byteValue, &cd); err != nil {
		return nil, err
	}

	if cd.Paths == nil {
		cd.Paths = make(map[string]WorkspaceData)
	}

	if cd.VersionAvailable == nil {
		cd.VersionAvailable = make(map[string]*Release)
	}

	return &cd, nil
}

func initFile(filepath string) error {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return os.WriteFile(filepath, []byte("{}"), 0644)
	}
	return err
}

func cacheSave(cd *CacheData) error {
	cachePath, err := getCachePath()
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(cd, "", "\t")
	if err != nil {
		return nil
	}

	return os.WriteFile(cachePath, bytes, 0644)
}
