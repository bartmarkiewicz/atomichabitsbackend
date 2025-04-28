package habit

import (
	"net/http"
)

type API struct{}

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
func (a *API) GetHabit(w http.ResponseWriter, r *http.Request) {}

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
func (a *API) CreateHabit(w http.ResponseWriter, r *http.Request) {}

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
func (a *API) GetHabits(w http.ResponseWriter, r *http.Request) {}

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
func (a *API) UpdateHabit(w http.ResponseWriter, r *http.Request) {}

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
func (a *API) DeleteHabit(w http.ResponseWriter, r *http.Request) {}
