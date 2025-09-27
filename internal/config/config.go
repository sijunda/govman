package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	viper "github.com/spf13/viper"
)

type Config struct {
	InstallDir     string           `mapstructure:"install_dir"`
	CacheDir       string           `mapstructure:"cache_dir"`
	DefaultVersion string           `mapstructure:"default_version"`
	Download       DownloadConfig   `mapstructure:"download"`
	Mirror         MirrorConfig     `mapstructure:"mirror"`
	AutoSwitch     AutoSwitchConfig `mapstructure:"auto_switch"`
	Shell          ShellConfig      `mapstructure:"shell"`
	GoReleases     GoReleasesConfig `mapstructure:"go_releases"`
	SelfUpdate     SelfUpdateConfig `mapstructure:"self_update"`
	Quiet          bool             `mapstructure:"quiet"`
	Verbose        bool             `mapstructure:"verbose"`
	configPath     string
}

type DownloadConfig struct {
	Parallel       bool          `mapstructure:"parallel"`
	MaxConnections int           `mapstructure:"max_connections"`
	Timeout        time.Duration `mapstructure:"timeout"`
	RetryCount     int           `mapstructure:"retry_count"`
	RetryDelay     time.Duration `mapstructure:"retry_delay"`
}

type MirrorConfig struct {
	Enabled bool   `mapstructure:"enabled"`
	URL     string `mapstructure:"url"`
}

type AutoSwitchConfig struct {
	Enabled     bool   `mapstructure:"enabled"`
	ProjectFile string `mapstructure:"project_file"`
}

type ShellConfig struct {
	AutoDetect bool `mapstructure:"auto_detect"`
	Completion bool `mapstructure:"completion"`
}

type GoReleasesConfig struct {
	APIURL      string        `mapstructure:"api_url"`
	DownloadURL string        `mapstructure:"download_url"`
	CacheExpiry time.Duration `mapstructure:"cache_expiry"`
}

type SelfUpdateConfig struct {
	GitHubAPIURL      string `mapstructure:"github_api_url"`
	GitHubReleasesURL string `mapstructure:"github_releases_url"`
}

func Load(configFile string) (*Config, error) {
	cfg := &Config{}

	// Set default values
	cfg.setDefaults()

	// Determine config file path
	if configFile != "" {
		cfg.configPath = configFile
	} else {
		homeDir, err := getHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		cfg.configPath = filepath.Join(homeDir, ".govman", "config.yaml")
	}

	// Setup viper
	viper.SetConfigFile(cfg.configPath)
	viper.SetConfigType("yaml")

	// Check if config file exists
	if _, err := os.Stat(cfg.configPath); os.IsNotExist(err) {
		// Config file doesn't exist, create it with default values
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("failed to create config file with default values: %w", err)
		}
	}

	// Read config file (it should now exist)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal config
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Expand paths
	if err := cfg.expandPaths(); err != nil {
		return nil, fmt.Errorf("failed to expand paths: %w", err)
	}

	// Create directories
	if err := cfg.createDirectories(); err != nil {
		return nil, fmt.Errorf("failed to create directories: %w", err)
	}

	return cfg, nil
}

func (c *Config) setDefaults() {
	homeDir, err := getHomeDir()
	if err != nil {
		homeDir = "." // Use current directory as fallback instead of /tmp
	}
	govmanDir := filepath.Join(homeDir, ".govman")

	c.InstallDir = filepath.Join(govmanDir, "versions")
	c.CacheDir = filepath.Join(govmanDir, "cache")
	c.DefaultVersion = ""
	c.Quiet = false
	c.Verbose = false

	c.Download = DownloadConfig{
		Parallel:       true,
		MaxConnections: 4,
		Timeout:        300 * time.Second,
		RetryCount:     3,
		RetryDelay:     5 * time.Second,
	}

	c.Mirror = MirrorConfig{
		Enabled: false,
		URL:     "https://golang.google.cn/dl/",
	}

	c.AutoSwitch = AutoSwitchConfig{
		Enabled:     true,
		ProjectFile: ".govman-version",
	}

	c.Shell = ShellConfig{
		AutoDetect: true,
		Completion: true,
	}

	c.GoReleases = GoReleasesConfig{
		APIURL:      "https://go.dev/dl/?mode=json&include=all",
		DownloadURL: "https://go.dev/dl/%s",
		CacheExpiry: 10 * time.Minute,
	}

	c.SelfUpdate = SelfUpdateConfig{
		GitHubAPIURL:      "https://api.github.com/repos/sijunda/govman/releases/latest",
		GitHubReleasesURL: "https://api.github.com/repos/sijunda/govman/releases?per_page=1",
	}
}

func (c *Config) expandPaths() error {
	var err error

	c.InstallDir, err = expandPath(c.InstallDir)
	if err != nil {
		return fmt.Errorf("failed to expand install_dir: %w", err)
	}

	c.CacheDir, err = expandPath(c.CacheDir)
	if err != nil {
		return fmt.Errorf("failed to expand cache_dir: %w", err)
	}

	return nil
}

func (c *Config) createDirectories() error {
	dirs := []string{c.InstallDir, c.CacheDir}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func (c *Config) Save() error {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(c.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Sync struct fields with viper before saving
	viper.Set("default_version", c.DefaultVersion)
	viper.Set("install_dir", c.InstallDir)
	viper.Set("cache_dir", c.CacheDir)
	viper.Set("quiet", c.Quiet)
	viper.Set("verbose", c.Verbose)
	viper.Set("download", c.Download)
	viper.Set("mirror", c.Mirror)
	viper.Set("auto_switch", c.AutoSwitch)
	viper.Set("shell", c.Shell)
	viper.Set("go_releases", c.GoReleases)
	viper.Set("self_update", c.SelfUpdate)

	// Write config file
	if err := viper.WriteConfigAs(c.configPath); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func (c *Config) GetVersionDir(version string) string {
	return filepath.Join(c.InstallDir, fmt.Sprintf("go%s", version))
}

func (c *Config) GetBinPath() string {
	homeDir, err := getHomeDir()
	if err != nil {
		// Fallback to current directory if home directory cannot be determined
		homeDir = "."
	}
	return filepath.Join(homeDir, ".govman", "bin")
}

func (c *Config) GetCurrentSymlink() string {
	return filepath.Join(c.GetBinPath(), "go")
}

func getHomeDir() (string, error) {
	var homeDir string
	if runtime.GOOS == "windows" {
		homeDir = os.Getenv("USERPROFILE")
	} else {
		homeDir = os.Getenv("HOME")
	}

	if homeDir == "" {
		return "", fmt.Errorf("unable to determine home directory: HOME/USERPROFILE environment variable is not set")
	}

	return homeDir, nil
}

func expandPath(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("empty path provided")
	}
	if path[0] == '~' {
		homeDir, err := getHomeDir()
		if err != nil {
			return "", err
		}
		// Validate that the path after ~ doesn't contain directory traversal
		if len(path) > 1 && (path[1] != '/' && path[1] != '\\') {
			return "", fmt.Errorf("invalid path format: paths starting with ~ must be followed by / or \\")
		}

		expandedPath := filepath.Join(homeDir, path[1:])
		// Ensure the expanded path is within the home directory to prevent directory traversal
		rel, err := filepath.Rel(homeDir, expandedPath)
		if err != nil {
			return "", fmt.Errorf("failed to evaluate relative path: %w", err)
		}
		if strings.HasPrefix(rel, "..") || strings.HasPrefix(filepath.ToSlash(rel), "../") {
			return "", fmt.Errorf("path traversal detected: expanded path is outside home directory")
		}

		return expandedPath, nil
	}
	return path, nil
}
