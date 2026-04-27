package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Storage  StorageConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Log      LogConfig
	Upload   UploadConfig
	App      AppConfig
	LLM      LLMConfig
	Lark     LarkConfig
}

type ServerConfig struct {
	Port int
	Mode string
}

type StorageConfig struct {
	Mode string // memory 或 database
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int // seconds
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

type JWTConfig struct {
	Secret      string
	ExpireHours int `mapstructure:"expire_hours"`
	Issuer      string
}

func (j JWTConfig) ExpireDuration() time.Duration {
	return time.Duration(j.ExpireHours) * time.Hour
}

type LogConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

type UploadConfig struct {
	MaxSize      int64
	AllowedTypes []string
	SavePath     string
}

type AppConfig struct {
	BaseURL              string
	QRCodeBaseURL        string
	GPSValidationEnabled bool
	GPSToleranceMeters   int
}

type LLMConfig struct {
	Provider string // openai, deepseek, ollama
	BaseURL  string
	APIKey   string
	Model    string
}

type LarkConfig struct {
	AppID             string `mapstructure:"app_id"`
	AppSecret         string `mapstructure:"app_secret"`
	VerificationToken string `mapstructure:"verification_token"`
	EncryptKey        string `mapstructure:"encrypt_key"`
}

var Cfg *Config

func Load(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	Cfg = &Config{}
	if err := viper.Unmarshal(Cfg); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := applyEnvOverrides(Cfg); err != nil {
		return err
	}

	return nil
}

func applyEnvOverrides(cfg *Config) error {
	overrideString(&cfg.Server.Mode, "EMS_SERVER_MODE", "SERVER_MODE", "GIN_MODE")
	overrideString(&cfg.Storage.Mode, "EMS_STORAGE_MODE", "STORAGE_MODE")

	if err := overrideInt(&cfg.Server.Port, "EMS_SERVER_PORT", "SERVER_PORT", "PORT"); err != nil {
		return err
	}

	overrideString(&cfg.Database.Host, "EMS_DATABASE_HOST", "DB_HOST")
	overrideString(&cfg.Database.User, "EMS_DATABASE_USER", "DB_USER")
	overrideString(&cfg.Database.Password, "EMS_DATABASE_PASSWORD", "DB_PASSWORD")
	overrideString(&cfg.Database.DBName, "EMS_DATABASE_NAME", "DB_NAME")
	overrideString(&cfg.Database.SSLMode, "EMS_DATABASE_SSLMODE", "DB_SSLMODE")

	if err := overrideInt(&cfg.Database.Port, "EMS_DATABASE_PORT", "DB_PORT"); err != nil {
		return err
	}
	if err := overrideInt(&cfg.Database.MaxIdleConns, "EMS_DATABASE_MAX_IDLE_CONNS"); err != nil {
		return err
	}
	if err := overrideInt(&cfg.Database.MaxOpenConns, "EMS_DATABASE_MAX_OPEN_CONNS"); err != nil {
		return err
	}
	if err := overrideInt(&cfg.Database.ConnMaxLifetime, "EMS_DATABASE_CONN_MAX_LIFETIME"); err != nil {
		return err
	}

	overrideString(&cfg.Redis.Host, "EMS_REDIS_HOST", "REDIS_HOST")
	overrideString(&cfg.Redis.Password, "EMS_REDIS_PASSWORD", "REDIS_PASSWORD")

	if err := overrideInt(&cfg.Redis.Port, "EMS_REDIS_PORT", "REDIS_PORT"); err != nil {
		return err
	}
	if err := overrideInt(&cfg.Redis.DB, "EMS_REDIS_DB", "REDIS_DB"); err != nil {
		return err
	}
	if err := overrideInt(&cfg.Redis.PoolSize, "EMS_REDIS_POOL_SIZE", "REDIS_POOL_SIZE"); err != nil {
		return err
	}

	overrideString(&cfg.JWT.Secret, "EMS_JWT_SECRET", "JWT_SECRET")
	overrideString(&cfg.JWT.Issuer, "EMS_JWT_ISSUER", "JWT_ISSUER")
	if err := overrideInt(&cfg.JWT.ExpireHours, "EMS_JWT_EXPIRE_HOURS", "JWT_EXPIRE_HOURS"); err != nil {
		return err
	}

	overrideString(&cfg.Log.Level, "EMS_LOG_LEVEL", "LOG_LEVEL")
	overrideString(&cfg.Log.Filename, "EMS_LOG_FILENAME", "LOG_FILENAME")
	if err := overrideInt(&cfg.Log.MaxSize, "EMS_LOG_MAX_SIZE"); err != nil {
		return err
	}
	if err := overrideInt(&cfg.Log.MaxBackups, "EMS_LOG_MAX_BACKUPS"); err != nil {
		return err
	}
	if err := overrideInt(&cfg.Log.MaxAge, "EMS_LOG_MAX_AGE"); err != nil {
		return err
	}
	if err := overrideBool(&cfg.Log.Compress, "EMS_LOG_COMPRESS"); err != nil {
		return err
	}

	if err := overrideInt64(&cfg.Upload.MaxSize, "EMS_UPLOAD_MAX_SIZE"); err != nil {
		return err
	}
	overrideString(&cfg.Upload.SavePath, "EMS_UPLOAD_SAVE_PATH", "UPLOAD_SAVE_PATH")

	overrideString(&cfg.App.BaseURL, "EMS_APP_BASE_URL", "APP_BASE_URL")
	overrideString(&cfg.App.QRCodeBaseURL, "EMS_APP_QR_CODE_BASE_URL", "QR_CODE_BASE_URL")
	if err := overrideBool(&cfg.App.GPSValidationEnabled, "EMS_APP_GPS_VALIDATION_ENABLED"); err != nil {
		return err
	}
	if err := overrideInt(&cfg.App.GPSToleranceMeters, "EMS_APP_GPS_TOLERANCE_METERS"); err != nil {
		return err
	}

	overrideString(&cfg.LLM.Provider, "EMS_LLM_PROVIDER", "LLM_PROVIDER")
	overrideString(&cfg.LLM.BaseURL, "EMS_LLM_BASE_URL", "LLM_BASE_URL")
	overrideString(&cfg.LLM.APIKey, "EMS_LLM_API_KEY", "LLM_API_KEY")
	overrideString(&cfg.LLM.Model, "EMS_LLM_MODEL", "LLM_MODEL")

	overrideString(&cfg.Lark.AppID, "EMS_LARK_APP_ID", "LARK_APP_ID")
	overrideString(&cfg.Lark.AppSecret, "EMS_LARK_APP_SECRET", "LARK_APP_SECRET")
	overrideString(&cfg.Lark.VerificationToken, "EMS_LARK_VERIFICATION_TOKEN", "LARK_VERIFICATION_TOKEN")
	overrideString(&cfg.Lark.EncryptKey, "EMS_LARK_ENCRYPT_KEY", "LARK_ENCRYPT_KEY")

	return nil
}

func overrideString(target *string, keys ...string) {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			*target = value
			return
		}
	}
}

func overrideInt(target *int, keys ...string) error {
	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		if !ok || value == "" {
			continue
		}

		parsed, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("invalid integer value for %s: %w", key, err)
		}

		*target = parsed
		return nil
	}

	return nil
}

func overrideInt64(target *int64, keys ...string) error {
	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		if !ok || value == "" {
			continue
		}

		parsed, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid int64 value for %s: %w", key, err)
		}

		*target = parsed
		return nil
	}

	return nil
}

func overrideBool(target *bool, keys ...string) error {
	for _, key := range keys {
		value, ok := os.LookupEnv(key)
		if !ok || value == "" {
			continue
		}

		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("invalid boolean value for %s: %w", key, err)
		}

		*target = parsed
		return nil
	}

	return nil
}
