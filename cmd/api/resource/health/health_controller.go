package health

import "net/http"

// HealthCheck Read godoc
//
//	@summary		Health check
//	@description	Health check
//	@tags			health
//	@success		200
//	@router			/health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {}
