package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"habitgobackend/cmd/api/resource/habit"
	"habitgobackend/cmd/api/resource/health"
)

func New(database *gorm.DB, validator *validator.Validate) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/health", health.HealthCheck)

	router.Route("/v1", func(router chi.Router) {
		habitAPI := habit.New(database, validator)
		router.Get("/habits", habitAPI.GetHabits)
		router.Post("/habits", habitAPI.CreateHabit)
		router.Get("/habits/{id}", habitAPI.GetHabit)
		router.Put("/habits/{id}", habitAPI.UpdateHabit)
		router.Delete("/habits/{id}", habitAPI.DeleteHabit)
	})

	return router
}
