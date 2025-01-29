package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string        `yaml:"env" env-default:"local"`
	StoragePath    string        `yaml:"storage_path" env-required:"true"`
	GRPC           GRPCConfig    `yaml:"grpc"`
	MigrationsPath string        `yaml:"migrations_path" env-default:"./migrations"`
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port" env-default:"44044"`
	Timeout time.Duration `yaml:"timeout" env-default:"5s"`
}

// MustLoad загружает конфигурацию и завершает программу в случае ошибки.
func MustLoad() *Config {
	configPath := fetchConfigPath()
	cfg, err := LoadByPath(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

// LoadByPath загружает конфигурацию по указанному пути.
func LoadByPath(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Загружаем конфигурацию окружения, если она указана
	if cfg.Env != "" {
		envConfigPath := fmt.Sprintf("%s.%s.yaml", "config", cfg.Env)
		if _, err := os.Stat(envConfigPath); err == nil {
			if err := cleanenv.ReadConfig(envConfigPath, &cfg); err != nil {
				return nil, fmt.Errorf("failed to read env config: %w", err)
			}
		}
	}

	// Валидация конфигурации
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validateConfig выполняет кастомную валидацию конфигурации.
func validateConfig(cfg *Config) error {
	if cfg.StoragePath == "" {
		return fmt.Errorf("storage_path is required")
	}
	if cfg.GRPC.Port == 0 {
		return fmt.Errorf("grpc.port is required")
	}
	if cfg.GRPC.Timeout == 0 {
		return fmt.Errorf("grpc.timeout is required")
	}
	return nil
}

// fetchConfigPath возвращает путь к конфигурационному файлу.
// Приоритет: флаг > переменная окружения > значение по умолчанию.
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
