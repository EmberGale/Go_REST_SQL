package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
)

func Load() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		if err := gotenv.Load(".env"); err != nil {
			return nil, fmt.Errorf("loading .env: %w", err)
		}
	}

	v := viper.New()
	v.AutomaticEnv()
	if err := bindEnvs(v, "", reflect.TypeFor[Config]()); err != nil {
		return nil, fmt.Errorf("binding env: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
	)); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	fmt.Printf("%#v\n", cfg)

	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func bindEnvs(v *viper.Viper, prefix string, typ reflect.Type) error {
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		key := field.Tag.Get("mapstructure")
		if key == "" {
			continue
		}
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}
		if field.Type.Kind() == reflect.Struct {
			if err := bindEnvs(v, fullKey, field.Type); err != nil {
				return err
			}
			continue
		}
		// viper doesn't reliably apply our nested-key formatting when binding by key,
		// so we bind to the exact env var name we expect, e.g.:
		//   fullKey=database.host -> APP_DATABASE__HOST
		envVar := "APP_" + strings.ToUpper(strings.ReplaceAll(fullKey, ".", "__"))
		if err := v.BindEnv(fullKey, envVar); err != nil {
			return fmt.Errorf("bind %s: %w", fullKey, err)
		}
	}
	return nil
}

func validateConfig(cfg *Config) error {
	if cfg.Database.Host == "" {
		return fmt.Errorf("required env variable cfg.Database.Host is not set")
	}
	if cfg.Database.Port == 0 {
		return fmt.Errorf("required env variable cfg.Database.Port is not set")
	}
	if cfg.Database.User == "" {
		return fmt.Errorf("required env variable cfg.Database.User is not set")
	}
	if cfg.Database.Password == "" {
		return fmt.Errorf("required env variable cfg.Database.Password is not set")
	}
	if cfg.Database.DBName == "" {
		return fmt.Errorf("required env variable cfg.Database.DBName is not set")
	}
	if cfg.Database.MigrationsPath == "" {
		return fmt.Errorf("required env variable cfg.Database.MigrationsPath is not set")
	}
	if cfg.Server.Port == "" {
		return fmt.Errorf("required env variable cfg.Server.Port is not set")
	}
	if cfg.Logger.Level == "" {
		return fmt.Errorf("required env variable cfg.Logger.Level is not set")
	}
	return nil
}
