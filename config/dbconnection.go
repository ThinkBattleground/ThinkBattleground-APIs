package config

import (
	"database/sql"
	"log"
	"os"
	"thinkbattleground-apis/constants"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func DbConnection() {
	var err error
	if err := LoadEnv(); err != nil {
		log.Println(constants.LOAD_ENV_ERROR)
		return
	}

	url := os.Getenv("URL")
	DB, err = sql.Open("postgres", url)
	if err != nil {
		log.Println("Error While DB Connection: ", err)
		panic(err)
	}
	if err = DB.Ping(); err != nil {
		log.Println("Error While Ping DB: ", err)
		panic(err)
	}

	log.Println("The database connected Successfully!")
}
