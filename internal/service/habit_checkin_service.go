package service

import (
	"errors"
	"habit-tracker-api/internal/domain"
	"habit-tracker-api/internal/repository"
	"time"
)

type HabitCheckinService struct {
	habitRepo   *repository.HabitRepository
	checkinRepo *repository.HabitCheckinRepository
}

func NewHabitCheckinService(
	hr *repository.HabitRepository,
	cr *repository.HabitCheckinRepository,
) *HabitCheckinService {
	return &HabitCheckinService{hr, cr}
}

// CheckIn отмечает выполнение привычки на сегодня (или в указанную дату)
func (s *HabitCheckinService) CheckIn(habitID, comment string) error {
	// Проверяем существование привычки
	if _, err := s.habitRepo.FindByID(habitID); err != nil {
		return errors.New("habit not found")
	}
	hc := &domain.HabitCheckin{
		HabitID: habitID,
		// Truncate обнуляет время, оставляя только дату
		Date:    time.Now().Truncate(24 * time.Hour),
		Comment: comment,
	}
	return s.checkinRepo.Create(hc)
}

// Stats возвращает статистику по привычке:
// streak — дней подряд до сегодня,
// totalChecks — уникальных дней выполнения,
// possibleChecks — дней с момента создания,
// completionRate — процент выполнения
func (s *HabitCheckinService) Stats(habitID string) (
	streak int, totalChecks, possibleChecks int, completionRate float64, err error,
) {
	// Получаем все отметки
	checks, err := s.checkinRepo.FindByHabit(habitID)
	if err != nil {
		return
	}
	// Узнаём дату создания привычки
	h, err := s.habitRepo.FindByID(habitID)
	if err != nil {
		return
	}
	start := h.CreatedAt.Truncate(24 * time.Hour)
	today := time.Now().Truncate(24 * time.Hour)

	// Считаем, сколько дней всего прошло
	days := int(today.Sub(start).Hours()/24) + 1
	possibleChecks = days

	// Собираем уникальные даты выполнения
	done := make(map[string]bool, len(checks))
	for _, c := range checks {
		key := c.Date.Format("2006-01-02")
		done[key] = true
	}
	totalChecks = len(done)

	// Считаем текущий streak
	streak = 0
	for d := today; !d.Before(start); d = d.AddDate(0, 0, -1) {
		if done[d.Format("2006-01-02")] {
			streak++
		} else {
			break
		}
	}

	// Процент выполнения
	completionRate = float64(totalChecks) / float64(possibleChecks) * 100
	return
}

// Добавление Report

type HabitReport struct {
	Streak         int     `json:"streak"`
	TotalChecks    int     `json:"total_checks"`
	PossibleChecks int     `json:"possible_checks"`
	CompletionRate float64 `json:"completion_rate"`
}

func (s *HabitCheckinService) Report(habitID string) (*HabitReport, error) {
	streak, total, possible, rate, err := s.Stats(habitID)
	if err != nil {
		return nil, err
	}
	return &HabitReport{
		Streak:         streak,
		TotalChecks:    total,
		PossibleChecks: possible,
		CompletionRate: rate,
	}, nil
}
