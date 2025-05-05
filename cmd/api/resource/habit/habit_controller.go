package habit

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

import "gorm.io/gorm"

type Api struct {
	repository *Repository
}

func New(db *gorm.DB) *Api {
	return &Api{
		repository: NewRepository(db),
	}
}

// GetHabit godoc
//
//	@summary		Get single habit
//	@description	Get habit by ID
//	@tags			habits
//	@accept			json
//	@produce		json
//	@success		200	{object}	JsonHabit
//	@failure		404 {object}	error.Error
//	@failure		500	{object}	error.Error
//	@router			/habits/{id} [get]
func (a *Api) GetHabit(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	habit, err := a.repository.GetHabit(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if habit == nil {
		http.Error(w, "Habit not found", http.StatusNotFound)
		return
	}

	jsonHabit := habit.ToJson()

	if err := json.NewEncoder(w).Encode(jsonHabit); err != nil {
		return
	}
}

// CreateHabit godoc
//
//	@summary		Create Habit
//	@description	Create Habit
//	@tags			habits
//	@accept			json
//	@produce		json
//	@param			body	body	JsonHabit	true	"JsonHabit"
//	@success		201
//	@failure		400	{object}	error.Error
//	@failure		422	{object}	error.Errors
//	@failure		500	{object}	error.Error
//	@router			/habits [post]
func (a *Api) CreateHabit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Trying to create habit")
	jsonHabit := &JsonHabit{}
	if err := json.NewDecoder(r.Body).Decode(jsonHabit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newHabit := jsonHabit.ToHabit()
	newHabit.ID = uuid.New()

	_, err := a.repository.CreateHabit(newHabit)

	if err != nil {
		return
	}

	fmt.Println("Created habit")

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", "/habits/"+newHabit.ID.String())
}

// GetHabits godoc
//
//	@summary		List habits
//	@description	List habits
//	@tags			habits
//	@accept			json
//	@produce		json
//	@success		200	{object}	JsonHabits
//	@failure		500	{object}	error.Error
//	@router			/habits [get]
func (a *Api) GetHabits(w http.ResponseWriter, r *http.Request) {
	habits, err := a.repository.GetHabits()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(habits) == 0 {
		_, err := fmt.Fprint(w, "[]")
		if err != nil {
			return
		}
	}

	if err := json.NewEncoder(w).Encode(habits); err != nil {
		return
	}
	w.WriteHeader(http.StatusOK)
}

// UpdateHabit godoc
//
//	@summary		Update habit
//	@description	Update habit
//	@tags			habits
//	@accept			json
//	@produce		json
//	@param			id		path	string		true	"Habit ID"
//	@param			body	body	JsonHabit	true	"JsonHabit"
//	@success		200
//	@failure		400	{object}	error.Error
//	@failure		404
//	@failure		422	{object}	error.Errors
//	@failure		500	{object}	error.Error
//	@router			/habits/{id} [put]
func (a *Api) UpdateHabit(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonHabit := &JsonHabit{}
	if err := json.NewDecoder(r.Body).Decode(jsonHabit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	habit := jsonHabit.ToHabit()
	habit.ID = id

	rows, err := a.repository.UpdateHabit(habit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "Habit not found", http.StatusNotFound)
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// DeleteHabit godoc
//
//	@summary		Delete habit
//	@description	Delete habit
//	@tags			habits
//	@accept			json
//	@produce		json
//	@param			id	path	string	true	"Habit ID"
//	@success		200
//	@failure		400	{object}	error.Error
//	@failure		404
//	@failure		500	{object}	error.Error
//	@router			/habits/{id} [delete]
func (a *Api) DeleteHabit(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		// handle later
		return
	}

	rows, err := a.repository.DeleteHabit(id)
	if err != nil {
		// handle later
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}
