package BGServer

import (
	"database/sql"
	"encoding/json"
	"gopkg.in/validator.v2"
	"log"
	"os"
)

type Config struct {
	DSN string `json:"DSN,omitempty" validate:"nonzero"`
}

var (
	cfg Config
	db *sql.DB
)

func initConfig(configFile string) error{
	file, err := os.Open(configFile)
	if err != nil{
		return err
	}
	decoder := json.NewDecoder(file)
	defer file.Close()
	cfg = Config{}
	if err = decoder.Decode(&cfg); err != nil{
		return err
	}
	if err = validator.Validate(cfg); err != nil {
		return err
	}
	return nil
}

func connect() *sql.DB {
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
