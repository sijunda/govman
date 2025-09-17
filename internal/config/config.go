package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	configPath     string
}

type DownloadConfig struct {
	Parallel       bool          `mapstructure:"parallel"`
	MaxConnections int           `mapstructure:"max_connections"`
	Timeout        time.Duration `mapstructure:"timeout"`
	RetryCount     int           `mapstructure:"retry_count"`
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

	// Read config file if it exists
	if _, err := os.Stat(cfg.configPath); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
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
	homeDir, _ := getHomeDir()
	govmanDir := filepath.Join(homeDir, ".govman")

	c.InstallDir = filepath.Join(govmanDir, "versions")
	c.CacheDir = filepath.Join(govmanDir, "cache")
	c.DefaultVersion = ""

	c.Download = DownloadConfig{
		Parallel:       true,
		MaxConnections: 4,
		Timeout:        300 * time.Second,
		RetryCount:     3,
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
	homeDir, _ := getHomeDir()
	return filepath.Join(homeDir, ".govman", "bin")
}

func (c *Config) GetCurrentSymlink() string {
	return filepath.Join(c.GetBinPath(), "go")
}

func getHomeDir() (string, error) {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE"), nil
	}
	return os.Getenv("HOME"), nil
}

func expandPath(path string) (string, error) {
	if path[0] == '~' {
		homeDir, err := getHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(homeDir, path[1:]), nil
	}
	return path, nil
}
