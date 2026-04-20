package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Storage  StorageConfig  `mapstructure:"storage"`
	Log      LogConfig      `mapstructure:"log"`
}

type ServerConfig struct {
	HTTPPort     int    `mapstructure:"http_port"`
	SessionSecret string `mapstructure:"session_secret"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type StorageConfig struct {
	UploadDir string `mapstructure:"upload_dir"`
}

type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

var GlobalConfig *Config

func Load(configPath string) (*Config, error) {
	v := viper.New()

	v.SetDefault("server.http_port", 8080)
	v.SetDefault("server.session_secret", "clipper-secret")
	v.SetDefault("database.path", "./data/clipper.db")
	v.SetDefault("storage.upload_dir", "./uploads")
	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "text")

	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	v.SetEnvPrefix("CLIPPER")
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	GlobalConfig = &cfg
	return &cfg, nil
}

func GetUploadDir() string {
	if GlobalConfig != nil {
		return GlobalConfig.Storage.UploadDir
	}
	return "./uploads"
}

func GetDataDir() string {
	if GlobalConfig != nil {
		return GlobalConfig.Database.Path
	}
	return "./data/clipper.db"
}