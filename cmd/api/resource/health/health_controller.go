package health

import "net/http"

// HealthCheckHandler Read godoc
//
//	@summary		Health check
//	@description	Health check
//	@tags			health
//	@success		200
//	@router			/health [get]
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
