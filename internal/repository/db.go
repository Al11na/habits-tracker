package repository

import (
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

const userBucket = "Users"

var DB *bolt.DB

// InitDB открывает или создаёт файл базы и инициализирует нужные «бадкеты»
func InitDB() {
	var err error
	DB, err = bolt.Open("habit_tracker.db", 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatalf("failed to open BoltDB: %v", err)
	}

	// Создадим бакет для хранения пользователей
	err = DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(userBucket))
		return err
	})
	if err != nil {
		log.Fatalf("failed to create bucket: %v", err)
	}
}
