package lib

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"manuel71sj/go-api-template/pkg/file"
)

var configPath = "config/config.yml"
var casbinModelPath = "config/casbin_model.conf"

var defaultConfig = Config{
	Name: "api-backend",
	Http: &HttpConfig{
		Host: "0.0.0.0",
		Port: 8080,
	},
	Log: &LogConfig{
		Level:       "debug",
		Directory:   "./logs",
		Development: true,
	},
	SuperAdmin: &SuperAdminConfig{},
	Auth:       &AuthConfig{},
	Casbin:     &CasbinConfig{Enable: false},
	Redis:      &RedisConfig{Host: "192.168.5.58", Port: 6379},
	Database: &DatabaseConfig{
		Parameters:   "charset=utf8mb4&parseTime=True&loc=Local&allowNativePasswords=true&timeout=5s",
		MaxLifetime:  7200,
		MaxOpenConns: 150,
		MaxIdleConns: 50,
	},
}

// Config Configuration are the available config value.
type Config struct {
	Name       string            `mapstructure:"Name"`
	Http       *HttpConfig       `mapstructure:"Http"`
	Log        *LogConfig        `mapstructure:"Log"`
	SuperAdmin *SuperAdminConfig `mapstructure:"SuperAdmin"`
	Auth       *AuthConfig       `mapstructure:"Auth"`
	Casbin     *CasbinConfig     `mapstructure:"Casbin"`
	Redis      *RedisConfig      `mapstructure:"Redis"`
	Database   *DatabaseConfig   `mapstructure:"Database"`
}

func NewConfig() Config {
	config := defaultConfig

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		panic(errors.Wrap(err, "Failed to read config"))
	}

	if err := viper.Unmarshal(&config); err != nil {
		panic(errors.Wrap(err, "Failed to unmarshal config"))
	}

	config.Casbin.Model = casbinModelPath

	return config
}

func SetConfigPath(path string) {
	if !file.IsFile(path) {
		panic("config filepath does not exist")
	}

	configPath = path
}

func SetConfigCasbinModelPath(path string) {
	if !file.IsFile(path) {
		panic("casbin model filepath does not exist")
	}

	casbinModelPath = path
}

type HttpConfig struct {
	Host string `mapstructure:"Host" validate:"ipv4"`
	Port int    `mapstructure:"Port" validate:"gte=1,lte=65535"`
}

func (c *HttpConfig) ListenAddr() string {
	if err := validator.New().Struct(c); err != nil {
		return "0.0.0.0:8080"
	}

	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// LogConfig
// LogLevel     : debug,info,warn,error,dpanic,panic,fatal : default info
// Format       : json, console : default json
// Directory    : Log storage path : default "./"
type LogConfig struct {
	Level       string `mapstructure:"Level"`
	Format      string `mapstructure:"Format"`
	Directory   string `mapstructure:"Directory"`
	Development bool   `mapstructure:"Development"`
}

type SuperAdminConfig struct {
	Username string `mapstructure:"Username"`
	RealName string `mapstructure:"RealName"`
	Password string `mapstructure:"Password"`
}

type CaptchaConfig struct {
	Enable     bool `mapstructure:"Enable"`
	Width      int  `mapstructure:"Width"`
	Height     int  `mapstructure:"Height"`
	NoiseCount int  `mapstructure:"NoiseCount"`
}

type AuthConfig struct {
	Enable             bool           `mapstructure:"Enable"`
	TokenExpired       int            `mapstructure:"TokenExpired"`
	IgnorePathPrefixes []string       `mapstructure:"IgnorePathPrefixes"`
	Captcha            *CaptchaConfig `mapstructure:"Captcha"`
}

type CasbinConfig struct {
	Enable             bool     `mapstructure:"Enable"`
	Debug              bool     `mapstructure:"Debug"`
	Model              string   `mapstructure:"Model"`
	AutoLoad           bool     `mapstructure:"AutoLoad"`
	AutoLoadInternal   int      `mapstructure:"AutoLoadInternal"`
	IgnorePathPrefixes []string `mapstructure:"IgnorePathPrefixes"`
}

type DatabaseConfig struct {
	Engine      string `mapstructure:"Engine"`
	Name        string `mapstructure:"Name"`
	Host        string `mapstructure:"Host"`
	Port        int    `mapstructure:"Port"`
	Username    string `mapstructure:"Username"`
	Password    string `mapstructure:"Password"`
	TablePrefix string `mapstructure:"TablePrefix"`
	Parameters  string `mapstructure:"Parameters"`

	MaxLifetime  int `mapstructure:"MaxLifetime"`
	MaxOpenConns int `mapstructure:"MaxOpenConns"`
	MaxIdleConns int `mapstructure:"MaxIdleConns"`
}

func (c DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s", c.Username, c.Password, c.Host, c.Port, c.Name, c.Parameters)
}

type RedisConfig struct {
	Host      string `mapstructure:"Host"`
	Port      int    `mapstructure:"Port"`
	Password  string `mapstructure:"Password"`
	KeyPrefix string `mapstructure:"KeyPrefix"`
}

func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
