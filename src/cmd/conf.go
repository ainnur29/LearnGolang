package main

import (
	"fmt"
	"golang-bulang-bolang/src/config"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/stretchr/testify/assert/yaml"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
)

type Config struct {
	Server    config.ServerOptions   `yaml:"server"`
	Logger    config.LoggerOptions   `yaml:"logger"`
	Postgres  config.DatabaseOptions `yaml:"postgres"`
	MySQL     config.DatabaseOptions `yaml:"mysql"`
	Redis     config.RedisOptions    `yaml:"redis"`
	Queries   config.QueriesOptions  `yaml:"queries"`
	Auth      config.AuthOptions     `yaml:"auth"`
	Scheduler SchedulerConfig        `yaml:"scheduler"`
}

type SchedulerConfig struct {
	Enabled bool                `yaml:"enabled"`
	Jobs    SchedulerJobsConfig `yaml:"jobs"`
}

type SchedulerJobsConfig struct {
	UserGenerator UserGeneratorJobConfig `yaml:"user_generator"`
}

type UserGeneratorJobConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Cron      string `yaml:"cron"`
	BatchSize int    `yaml:"batch_size"`
	MinAge    int    `yaml:"min_age"`
	MaxAge    int    `yaml:"max_age"`
}

func InitConfig() (*Config, error) {
	cfgPath := "config.yaml"

	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, err
	}

	fmt.Println("masuk init conf : ", err)
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	overrideWithEnv(&cfg)

	return &cfg, nil
}

func overrideWithEnv(cfg *Config) {
	if val := os.Getenv("SERVER_PORT"); val != "" {
		cfg.Server.Port = parseInt(val, cfg.Server.Port)
	}

	if val := os.Getenv("LOG_LEVEL"); val != "" {
		cfg.Logger.Level = val
	}

	if val := os.Getenv("POSTGRES_HOST"); val != "" {
		cfg.Postgres.Host = val
	}

	if val := os.Getenv("POSTGRES_PORT"); val != "" {
		cfg.Postgres.Port = parseInt(val, cfg.Postgres.Port)
	}

	if val := os.Getenv("POSTGRES_USER"); val != "" {
		cfg.Postgres.User = val
	}

	if val := os.Getenv("POSTGRES_PASSWORD"); val != "" {
		cfg.Postgres.Password = val
	}

	if val := os.Getenv("POSTGRES_DB_NAME"); val != "" {
		cfg.Postgres.DBName = val
	}
}

func parseInt(s string, defaultVal int) int {
	var val int
	if _, err := fmt.Sscanf(s, "%d", &val); err == nil {
		return val
	}

	return defaultVal
}
