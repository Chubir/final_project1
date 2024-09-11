package config

import "os"

type Config struct {
	Port   string
	Dbfile string
}

func Init() *Config {
	conf := Config{}

	conf.Port = os.Getenv("TODO_PORT")
	if conf.Port == "" {
		conf.Port = "7540"
	}

	conf.Dbfile = os.Getenv("TODO_DBFILE")
	if conf.Dbfile == "" {
		conf.Dbfile = "scheduler.db"
	}
	return &conf
}
