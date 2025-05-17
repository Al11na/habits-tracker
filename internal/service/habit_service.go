package service

import (
	"errors"
	"habit-tracker-api/internal/domain"
	"habit-tracker-api/internal/repository"
	"strings"
	"time"
)

type HabitService struct {
	repo *repository.HabitRepository
}

func NewHabitService(r *repository.HabitRepository) *HabitService {
	return &HabitService{r}
}

func (s *HabitService) Create(h *domain.Habit) error {
	if h.Name == "" {
		return errors.New("habit name is required")
	}
	return s.repo.Create(h)
}

func (s *HabitService) GetAll(userEmail string) ([]*domain.Habit, error) {
	return s.repo.FindAllByUser(userEmail)
}

func (s *HabitService) GetByID(id string) (*domain.Habit, error) {
	return s.repo.FindByID(id)
}

func (s *HabitService) Update(h *domain.Habit) error {
	if h.ID == "" {
		return errors.New("habit ID is required")
	}
	return s.repo.Update(h)
}

func (s *HabitService) Delete(id string) error {
	return s.repo.Delete(id)
}

// ListHabits — возвращает отфильтрованный и постраничный список привычек пользователя
func (s *HabitService) ListHabits(
	userEmail, name, dateFrom, dateTo string,
	page, pageSize int,
) ([]*domain.Habit, error) {
	all, err := s.repo.FindAllByUser(userEmail)
	if err != nil {
		return nil, err
	}

	// фильтрация по имени и диапазону дат
	var filtered []*domain.Habit
	for _, h := range all {
		if name != "" && !strings.Contains(
			strings.ToLower(h.Name),
			strings.ToLower(name),
		) {
			continue
		}
		if dateFrom != "" {
			from, err := time.Parse("2006-01-02", dateFrom)
			if err != nil {
				return nil, errors.New("invalid date_from")
			}
			if h.CreatedAt.Before(from) {
				continue
			}
		}
		if dateTo != "" {
			to, err := time.Parse("2006-01-02", dateTo)
			if err != nil {
				return nil, errors.New("invalid date_to")
			}
			// включаем всю дату dateTo
			if h.CreatedAt.After(to.Add(24 * time.Hour)) {
				continue
			}
		}
		filtered = append(filtered, h)
	}

	// пагинация
	total := len(filtered)
	start := (page - 1) * pageSize
	if start > total {
		return []*domain.Habit{}, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return filtered[start:end], nil
}
