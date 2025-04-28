package habit

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
