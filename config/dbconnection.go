package config

import (
	"database/sql"
	"log"
	"os"
	"thinkbattleground-apis/constants"
	"time"

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
	
	for i := 0; i < 10; i++ { // retry 10 times
		DB, err = sql.Open("postgres", url)
		if err != nil {
			log.Println("DB connection failed:", err)
		} else if err = DB.Ping(); err == nil {
			log.Println("The database connected Successfully!")
			return
		}

		log.Println("Retrying DB connection in 3 seconds...")
		time.Sleep(3 * time.Second)
	}

	log.Fatal("Could not connect to DB after several attempts:", err)
}
