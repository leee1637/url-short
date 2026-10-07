package config

import "fmt"

type Config struct {
	Env        string     `yaml:"env" env:"APP_ENV"`
	HTTPServer HTTPServer `yaml:"http_server"`
	Postgres   Postgres   `yaml:"postgres"`
	Auth       Auth       `yaml:"auth"`
}

type Auth struct {
	User string `yaml:"user" env:"AUTH_USER" env-required:"true"`
	Pass string `yaml:"pass" env:"AUTH_PASS" env-required:"true"`
}

type HTTPServer struct {
	Addr string `yaml:"addr" env:"HTTP_ADDR" env-default:":8080"`
}

type Postgres struct {
	Host string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"DB_PORT" env-default:"5432"`
	User string `yaml:"user" env:"DB_USER" env-required:"true"`
	Pass string `yaml:"pass" env:"DB_PASS" env-required:"true"`
	Name string `yaml:"name" env:"DB_NAME" env-required:"true"`
}

func (p Postgres) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		p.User, p.Pass, p.Host, p.Port, p.Name)
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig("config/local.yaml", &cfg); err != nil {
		panic(err)
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil { // env перезаписывает yaml
		panic(err)
	}
	return &cfg
}
