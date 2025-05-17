package repository

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"time"

	"habit-tracker-api/internal/domain"

	"github.com/google/uuid"
	bolt "go.etcd.io/bbolt"
)

const habitBucket = "habits"

// HabitRepository работает поверх BoltDB
type HabitRepository struct{}

// NewHabitRepository — конструктор, создает bucket, если нужно
func NewHabitRepository() *HabitRepository {
	_ = DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(habitBucket))
		return err
	})
	return &HabitRepository{}
}

// Create — сохраняет новую привычку в BoltDB
func (r *HabitRepository) Create(h *domain.Habit) error {
	if h.ID == "" {
		h.ID = uuid.New().String()
	}
	if h.CreatedAt.IsZero() {
		h.CreatedAt = time.Now()
	}
	return DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(habitBucket))
		data, err := json.Marshal(h)
		if err != nil {
			return err
		}
		return b.Put([]byte(h.ID), data)
	})
}

// FindAllByUser — возвращает все привычки заданного пользователя
func (r *HabitRepository) FindAllByUser(email string) ([]*domain.Habit, error) {
	var habits []*domain.Habit
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(habitBucket))
		return b.ForEach(func(_, v []byte) error {
			var h domain.Habit
			if err := json.Unmarshal(v, &h); err != nil {
				return err
			}
			if h.UserEmail == email {
				habits = append(habits, &h)
			}
			return nil
		})
	})
	return habits, err
}

// FindByID — возвращает привычку по её ID
func (r *HabitRepository) FindByID(id string) (*domain.Habit, error) {
	var h domain.Habit
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(habitBucket))
		v := b.Get([]byte(id))
		if v == nil {
			return errors.New("habit not found")
		}
		return json.Unmarshal(v, &h)
	})
	if err != nil {
		return nil, err
	}
	return &h, nil
}

// Update — обновляет существующую привычку
func (r *HabitRepository) Update(h *domain.Habit) error {
	return DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(habitBucket))
		if b.Get([]byte(h.ID)) == nil {
			return errors.New("habit not found")
		}
		data, err := json.Marshal(h)
		if err != nil {
			return err
		}
		return b.Put([]byte(h.ID), data)
	})
}

// Delete — удаляет привычку по ID
func (r *HabitRepository) Delete(id string) error {
	return DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(habitBucket))
		if b.Get([]byte(id)) == nil {
			return errors.New("habit not found")
		}
		return b.Delete([]byte(id))
	})
}

// FindAll — возвращает все привычки с фильтрацией и пагинацией
func (r *HabitRepository) FindAll(
	name string,
	dateFrom, dateTo *time.Time,
	offset, limit int,
) ([]*domain.Habit, error) {
	var all []*domain.Habit

	// Сначала читаем всё
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(habitBucket))
		return b.ForEach(func(_, v []byte) error {
			var h domain.Habit
			if err := json.Unmarshal(v, &h); err != nil {
				return err
			}
			// Фильтрация по имени
			if name != "" && !containsIgnoreCase(h.Name, name) {
				return nil
			}
			// Фильтрация по диапазону дат
			if dateFrom != nil && h.CreatedAt.Before(*dateFrom) {
				return nil
			}
			if dateTo != nil && h.CreatedAt.After(*dateTo) {
				return nil
			}
			all = append(all, &h)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	// Сортируем по CreatedAt (последние первыми)
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.After(all[j].CreatedAt)
	})

	// Пагинация
	start := offset
	if start > len(all) {
		start = len(all)
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}

	return all[start:end], nil
}

// containsIgnoreCase проверяет вхождение подстроки без учета регистра
func containsIgnoreCase(s, substr string) bool {
	return strings.Contains(
		strings.ToLower(s),
		strings.ToLower(substr),
	)
}
