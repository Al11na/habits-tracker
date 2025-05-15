package repository

import (
	"bytes"
	"encoding/json"

	"habit-tracker-api/internal/domain"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

const checkinBucket = "habit_checkins"

type HabitCheckinRepository struct{}

func NewHabitCheckinRepository() *HabitCheckinRepository {
	// Создаём бакет, если надо
	_ = DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(checkinBucket))
		return err
	})
	return &HabitCheckinRepository{}
}

// Create — добавляет новую запись о выполнении
func (r *HabitCheckinRepository) Create(hc *domain.HabitCheckin) error {
	hc.ID = uuid.New().String()
	data, err := json.Marshal(hc)
	if err != nil {
		return err
	}
	return DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(checkinBucket))
		// ключ — сочетание habitID + дата(YYYY-MM-DD), чтобы избежать дубликатов
		key := []byte(hc.HabitID + "|" + hc.Date.Format("2006-01-02"))
		return b.Put(key, data)
	})
}

// FindByHabit — возвращает все check‑in’ы для данной привычки
func (r *HabitCheckinRepository) FindByHabit(habitID string) ([]domain.HabitCheckin, error) {
	var res []domain.HabitCheckin
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(checkinBucket))
		return b.ForEach(func(k, v []byte) error {
			// ключ начинается с habitID|
			if !bytes.HasPrefix(k, []byte(habitID+"|")) {
				return nil
			}
			var hc domain.HabitCheckin
			if err := json.Unmarshal(v, &hc); err != nil {
				return err
			}
			res = append(res, hc)
			return nil
		})
	})
	return res, err
}
