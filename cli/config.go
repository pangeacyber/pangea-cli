package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/spf13/viper"
)

type Service struct {
	Domain string `mapstructure:"domain"`
	Token  string `mapstructure:"token"`
}

func getEmptyService() Service {
	return Service{
		Token: "",
	}
}

type Profile map[string]Service

type ConfigFile struct {
	Title    string             `mapstructure:"title"`
	Version  string             `mapstructure:"version"`
	Profile  string             `mapstructure:"profile"`
	Profiles map[string]Profile `mapstructure:"profiles"`
}

var defaultConfig = ConfigFile{
	Title:    "Pangea",
	Version:  "v2.0",
	Profile:  "default",
	Profiles: map[string]Profile{},
}

func GetDefaultPangeaConfig() pangea.Config {
	return pangea.Config{
		Enviroment:         "production",
		PollResultTimeout:  time.Second * 30,
		QueuedRetryEnabled: true,
		CustomUserAgent:    "pangea-cli/" + Version,
	}
}

var ErrNoConfigFile = errors.New("pangea Token doesn't exist. Run `pangea login` to setup your CLI")
var ErrUnauthorized = errors.New("unauthorized token")

func (cf *ConfigFile) toMap() (map[string]any, error) {
	var resp map[string]any
	b, err := json.Marshal(cf)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(b, &resp)
	return resp, nil
}

func (cf *ConfigFile) validate() {
	if len(cf.Profiles) == 0 {
		cf.Profile = "default"
	}

	keys := make([]string, 0, len(cf.Profiles))
	for k := range cf.Profiles {
		keys = append(keys, k)
	}

	if len(keys) == 1 {
		cf.Profile = keys[0]
	}

	if len(keys) > 1 { // Check if current Profile exists, if not, set last of the list as selected
		var k string
		for _, k = range keys {
			if k == cf.Profile {
				break
			}
		}
		// if break when they are equal, re-assign same value. If no profile match, set the lastone as default
		cf.Profile = k
	}
}

func (cf *ConfigFile) save() error {
	configPath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	cf.validate()

	v := viper.New()
	v.SetConfigFile(configPath)

	m, err := cf.toMap()
	if err != nil {
		return err
	}

	err = v.MergeConfigMap(m)
	if err != nil {
		return err
	}

	if err := v.WriteConfig(); err != nil {
		return err
	}

	_ = cf.saveYAML()

	return nil
}

func (cf *ConfigFile) saveYAML() error {
	configPath, err := getConfigFilePathYAML()
	if err != nil {
		return err
	}
	cf.validate()

	v := viper.New()
	v.SetConfigFile(configPath)

	m, err := cf.toMap()
	if err != nil {
		return err
	}

	err = v.MergeConfigMap(m)
	if err != nil {
		return err
	}

	if err := v.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func getConfigFileDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Create the .pangea folder if it doesn't exist
	pangeaDir := filepath.Join(homeDir, ".pangea")
	err = os.MkdirAll(pangeaDir, 0700)
	if err != nil {
		return "", err
	}
	return pangeaDir, nil
}

func getConfigFilePath() (string, error) {
	pangeaDir, err := getConfigFileDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(pangeaDir, "config.toml"), nil
}

func getConfigFilePathYAML() (string, error) {
	pangeaDir, err := getConfigFileDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(pangeaDir, "config.yaml"), nil
}

func checkFileAndCreateDefault(configPath string) {
	if _, err := os.Stat(configPath); err == nil {
		return
	}

	_ = defaultConfig.save()
}

func loadConfig() (*ConfigFile, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	checkFileAndCreateDefault(configPath)

	v := viper.New()
	v.SetConfigFile(configPath)
	v.SetDefault("profiles", make(map[string]Profile))
	// Check if the config file exists
	if _, err := os.Stat(configPath); err != nil {
		return nil, ErrNoConfigFile
	}

	// Read the configuration
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error fetching your pangea token: %s", err)
	}

	// Get the token from the configuration
	var config ConfigFile
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %s", err)
	}

	// if config.Profiles == nil {
	// 	config.Profiles = make(map[string]Profile, 0)
	// }

	return &config, nil
}

type ps struct {
	p string // Profile name
	s string // Service name
}

func getPriorityList(profile, service string) []ps {
	return []ps{
		{p: profile, s: service},
		{p: profile, s: "default"},
		{p: "default", s: service},
		{p: "default", s: "default"},
	}
}

func GetTokenAndDomain(service string) (string, string, error) {
	config, err := loadConfig()
	if err != nil {
		return "", "", err
	}

	token, domain, err := config.GetTokenAndDomain(config.Profile, service)
	if err != nil {
		return "", "", err
	}
	return token, domain, nil
}

func GetProfileTokenAndDomain(profile, service string) (string, string, error) {
	config, err := loadConfig()
	if err != nil {
		return "", "", err
	}

	token, domain, err := config.GetTokenAndDomain(profile, service)
	if err != nil {
		return "", "", err
	}
	return token, domain, nil
}

func (cf *ConfigFile) GetTokenAndDomain(profile, service string) (string, string, error) {
	priority := getPriorityList(profile, service)
	token := ""
	domain := ""

	for _, ps := range priority {
		p, ok := cf.Profiles[ps.p]
		if !ok {
			continue
		}

		s, ok := p[ps.s]
		if !ok {
			continue
		}

		if token == "" {
			token = s.Token
		}

		if domain == "" {
			domain = s.Domain
		}

		if token != "" && domain != "" {
			break
		}
	}

	if token == "" {
		token = os.Getenv("PANGEA_TOKEN")
	}

	if token == "" {
		return "", "", fmt.Errorf("invalid empty token on profile:'%s' service:'%s'", cf.Profile, service)
	}

	if domain == "" {
		domain = os.Getenv("PANGEA_DOMAIN")
	}

	if domain == "" {
		return "", "", fmt.Errorf("empty domain on profile:'%s' service:'%s'", cf.Profile, service)
	}

	return token, domain, nil
}

func (cf *ConfigFile) initService(profile, service string) (string, string) {
	if profile == "" {
		profile = cf.Profile
	}

	if service == "" {
		service = "default"
	}

	_, ok := cf.Profiles[profile]
	if !ok {
		cf.Profiles[profile] = Profile{}
	}

	_, ok = cf.Profiles[profile][service]
	if !ok {
		cf.Profiles[profile][service] = getEmptyService()
	}
	return profile, service
}

func (cf *ConfigFile) SetToken(profile, service, token string) {
	profile, service = cf.initService(profile, service)

	s := cf.Profiles[profile][service]
	s.Token = token
	cf.Profiles[profile][service] = s
	_ = cf.save()
}

func (cf *ConfigFile) SetDomain(profile, service, domain string) {
	profile, service = cf.initService(profile, service)

	s := cf.Profiles[profile][service]
	s.Domain = domain
	cf.Profiles[profile][service] = s
	_ = cf.save()
}

// writeTokenToFile writes the Pangea token to the specified file
func SaveToken(profile, service, token string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	config.SetToken(profile, service, token)
	return config.save()
}

func SaveDomain(profile, service, domain string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	config.SetDomain(profile, service, domain)
	return config.save()
}

func ListProfiles() ([]string, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}

	keys := make([]string, 0, len(config.Profiles))
	for k := range config.Profiles {
		keys = append(keys, k)
	}
	return keys, nil
}

func GetCurrentProfile() (*Profile, string, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, "", err
	}

	p, ok := config.Profiles[config.Profile]
	if !ok {
		return nil, "", fmt.Errorf("current profile no available")
	}

	return &p, config.Profile, nil
}

func GetCurrentProfileName() (string, error) {
	config, err := loadConfig()
	if err != nil {
		return "", err
	}

	return config.Profile, nil
}

func CreateProfile(profileName string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	if profileName == "" {
		return errors.New("invalid emtpy profile name")
	}

	if _, ok := config.Profiles[profileName]; ok {
		return fmt.Errorf("profile '%s' already exists", profileName)
	}

	config.Profiles[profileName] = Profile{
		"default": getEmptyService(),
	}

	return config.save()
}

func DeleteProfile(profileName string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	if profileName == "" {
		return errors.New("invalid emtpy profile name")
	}

	if _, ok := config.Profiles[profileName]; !ok {
		return fmt.Errorf("profile '%s' do not exists", profileName)
	}

	delete(config.Profiles, profileName)
	config.validate()
	return config.save()
}

func SelectProfile(profileName string) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}

	if profileName == "" {
		return errors.New("invalid empty profile name")
	}

	if _, ok := config.Profiles[profileName]; !ok {
		return fmt.Errorf("profile '%s' does not exist on config file", profileName)
	}

	config.Profile = profileName
	return config.save()
}
