package handlers

import (
	"net/http"

	"encoding/json"
	"oauth2-provider/settings"
	"oauth2-provider/models"
	"sync"
)

type HealthCheckHandler struct {
	*settings.Oauth2ProviderSettings
}

// TODO In the future we could report back on the status of our DB, or our cache
// TODO (e.g. Redis) by performing a simple PING, and include them in the response.
func (h *HealthCheckHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var head string
	head, req.URL.Path = ShiftPath(req.URL.Path)

	switch head {
	case "":
		switch req.Method {
		case "GET", "POST":
			healthCheckStatus := h.isHealthy()
			if bytes, err := json.Marshal(healthCheckStatus); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			} else {
				if healthCheckStatus.Alive.Healthy {
					w.WriteHeader(http.StatusOK)
				} else {
					w.WriteHeader(http.StatusInternalServerError)
				}
				w.Write(bytes)
			}
		default:
			http.Error(w, "Not found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *HealthCheckHandler) isHealthy() *models.HealthCheckStatus {

	//response
	status := models.NewHealthCheckStatus()

	var wg sync.WaitGroup
	wg.Add(3)

	//Check KeyValueStore access
	go func(wg *sync.WaitGroup) {
		if code := h.Code(&models.AuthorizationRequest{ClientId: "health_check"}); code != "" {
			status.KvsAccess = &models.Healthy{Healthy: true}
		}
		wg.Done()
	}(&wg)

	//Check ClientId access
	go func(wg *sync.WaitGroup) {
		if _, err := h.GetClientInformation("health_check"); err == nil {
			status.ClientAccess = &models.Healthy{Healthy: true}
		}
		wg.Done()
	}(&wg)

	//Check user access management
	go func(wg *sync.WaitGroup) {
		if _, ok := h.MatchingCredentials("health_check", "check_health"); ok {
			status.UserAccess = &models.Healthy{Healthy: true}
		}
		wg.Done()
	}(&wg)

	wg.Wait()

	//Global Health check status
	if status.ClientAccess.Healthy && status.UserAccess.Healthy && status.KvsAccess.Healthy {
		status.Alive = &models.Healthy{Healthy: true}
	}

	return status
}
