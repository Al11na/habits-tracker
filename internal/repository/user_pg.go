package repository

import (
	"encoding/json"
	"errors"

	"habit-tracker-api/internal/domain"

	bolt "go.etcd.io/bbolt"
)

type UserRepository struct{}

// NewUserRepository просто гарантирует, что DB инициализирована
func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(user *domain.User) error {
	return DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))
		// проверим, нет ли такого email
		if b.Get([]byte(user.Email)) != nil {
			return errors.New("user already exists")
		}
		data, _ := json.Marshal(user)
		return b.Put([]byte(user.Email), data)
	})
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(userBucket))
		v := b.Get([]byte(email))
		if v == nil {
			return errors.New("user not found")
		}
		return json.Unmarshal(v, &user)
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}
