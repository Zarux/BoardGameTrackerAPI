package server

import (
	"database/sql"
	"gopkg.in/validator.v2"
	"log"
	"os"
)

type Config struct {
	DSN string `json:"DSN,omitempty" validate:"nonzero"`
}

var (
	cfg Config
	db  *sql.DB
)

func GetConfig() Config {
	if cfg == (Config{}) {
		err := InitConfig()
		if err != nil {
			log.Println(err)
			return Config{}
		}
	}
	return cfg
}

func InitConfig() error {
	cfg = Config{}
	cfg.DSN = os.Getenv("DSN_MAIN")
	if cfg.DSN == "" {
		log.Fatal("missing environment variable DSN_MAIN")
	}

	if err := validator.Validate(cfg); err != nil {
		return err
	}
	return nil
}

func connect() *sql.DB {
	cfg = GetConfig()
	var err error
	if db != nil {
		return db
	}
	db, err = sql.Open("mysql", cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db
}
