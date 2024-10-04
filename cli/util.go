package cli

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TylerBrock/colorjson"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

// Release represents a release asset for current OS and arch.
type Release struct {
	// Version is the version of the release
	Version string `json:"version"`

	// AssetURL is a URL to the uploaded file for the release
	AssetURL string `json:"asset_url"`

	// URL is a URL to release page for browsing
	URL string `json:"url"`

	// ReleaseNotes is a release notes of the release
	ReleaseNotes string `json:"release_notes"`

	// Name represents a name of the release
	Name string `json:"name"`

	// PublishedAt is the time when the release was published
	PublishedAt string `json:"published_at"`
}

func NewRelease(r *selfupdate.Release) *Release {
	if r == nil {
		return nil
	}

	return &Release{
		Version:      r.Version.String(),
		AssetURL:     r.AssetURL,
		URL:          r.URL,
		ReleaseNotes: r.ReleaseNotes,
		Name:         r.Name,
		PublishedAt:  r.PublishedAt.Format(time.RFC3339),
	}
}

func PrettyPrint(obj any) {
	f := colorjson.NewFormatter()
	f.Indent = 2

	s, _ := f.Marshal(obj)
	fmt.Println(string(s))
}

func IndentedString(obj any) (string, error) {
	b, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func GetCacheFolder() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	cacheDir := filepath.Join(homeDir, ".pangea", "cache")
	return cacheDir, nil
}

func RemoveCacheFolder() error {
	folder, err := GetCacheFolder()
	if err != nil {
		return err
	}
	return os.RemoveAll(folder)
}

func ReadAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

// ReadStdin reads user input from the command line
func ReadStdin() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// Logger that omits timestamp and still allows us to log fatal
var logger = log.New(os.Stderr, "", 0)

func GetLogger() *log.Logger {
	return logger
}

func GetDaySinceEpoch() string {
	// Day since epoch
	return fmt.Sprint(time.Now().UnixMilli() / 1000 / 3600 / 24)
}
