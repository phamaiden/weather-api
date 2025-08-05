package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/phamaiden/weather-api/internal/services"
)

type WeatherHandler struct {
	weatherService services.WeatherService
}

func NewWeatherHandler(ws services.WeatherService) *WeatherHandler {
	return &WeatherHandler{
		weatherService: ws,
	}
}

func (h *WeatherHandler) GetWeatherByCity(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	city, ok := vars["city"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "missing city parameter",
		})
		return
	}

	weather, err := h.weatherService.GetWeatherByCity(r.Context(), city)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weather)
}
