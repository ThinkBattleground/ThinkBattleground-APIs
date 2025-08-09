package config

import (
	"fmt"
	"log"
	"sync/atomic"

	"github.com/robfig/cron/v3"
)

func CleanupOldRecords(table string) {
	// SQL query to delete records older than 2 minutes
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE created_at < NOW() - INTERVAL '2 minutes';
	`, table)

	result, err := DB.Exec(query)
	if err != nil {
		log.Printf("Failed to delete old records: %v", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	log.Printf("Deleted %d old records from temp_users table.", rowsAffected)
}

func CronSchedule(table string) {
	var retryCount int32

	c := cron.New()
	_, err := c.AddFunc("@every 2m", func() {
		currentRetry := atomic.AddInt32(&retryCount, 1)

		if currentRetry >= 2 {
			fmt.Println("Max retries reached. Stopping cron scheduler.")
			c.Stop()
		} else {
			fmt.Println("Performing cleanup...")
			CleanupOldRecords(table)
		}
	})
	if err != nil {
		log.Fatalf("Failed to schedule cleanup job: %v", err)
	}
	c.Start()

	// Keep the application running
	select {}
}
