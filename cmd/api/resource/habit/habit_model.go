package habit

import "github.com/google/uuid"

type JsonHabit struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	ColourHex   string `json:"colourHex"`
	IconBase64  string `json:"iconBase64"`
	ModeType    string `json:"modeType"`
}

type JsonHabits struct {
	Habits []JsonHabit `json:"habits"`
}

type Habit struct {
	ID          uuid.UUID `gorm:"primary_key"`
	Description string
	ColourHex   string
	IconBase64  string
	ModeType    string
}

type Habits []*Habit

func (h Habit) ToJson() JsonHabit {
	return JsonHabit{
		ID:          h.ID.String(),
		Description: h.Description,
		ColourHex:   h.ColourHex,
		IconBase64:  h.IconBase64,
		ModeType:    h.ModeType,
	}
}

func (h JsonHabit) ToHabit() *Habit {
	id, _ := uuid.Parse(h.ID)

	return &Habit{
		ID:          id,
		Description: h.Description,
		ColourHex:   h.ColourHex,
		IconBase64:  h.IconBase64,
		ModeType:    h.ModeType,
	}
}
