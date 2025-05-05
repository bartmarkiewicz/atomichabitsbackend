package habit

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	database *gorm.DB
}

func NewRepository(database *gorm.DB) *Repository {
	return &Repository{database}
}

func (repository *Repository) GetHabits() (Habits, error) {
	habits := make([]*Habit, 0)
	if err := repository.database.Find(&habits).Error; err != nil {
		return nil, err
	}
	return habits, nil
}

func (repository *Repository) CreateHabit(habit *Habit) (*Habit, error) {
	if err := repository.database.Create(habit).Error; err != nil {
		return nil, err
	}
	return habit, nil
}

func (repository *Repository) GetHabit(id uuid.UUID) (*Habit, error) {
	habit := &Habit{}
	if err := repository.database.
		Where("id = ?", id).
		First(&habit).Error; err != nil {
		return nil, err
	}
	return habit, nil
}

func (repository *Repository) UpdateHabit(habit *Habit) (int64, error) {
	result := repository.database.
		Model(&Habit{}).
		Select("Description", "ColourHex", "IconBase64", "ModeType").
		Where("id = ?", habit.ID).
		Updates(habit)

	return result.RowsAffected, result.Error
}

func (repository *Repository) DeleteHabit(id uuid.UUID) (int64, error) {
	result := repository.database.Where("id = ?", id).Delete(&Habit{})

	return result.RowsAffected, result.Error
}
