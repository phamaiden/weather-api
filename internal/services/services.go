package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/phamaiden/weather-api/internal/cache"
	"github.com/phamaiden/weather-api/internal/models"
)

type WeatherService interface {
	GetWeatherByCity(ctx context.Context, city string) (models.WeatherResponse, error)
}

type weatherService struct {
	apiKey string
	cache  cache.Cache
}

func NewWeatherService(key string, cache *cache.Cache) *weatherService {
	return &weatherService{
		apiKey: key,
		cache:  *cache,
	}
}

func (s *weatherService) GetWeatherByCity(ctx context.Context, city string) (models.WeatherResponse, error) {
	// Check Cache
	cacheKey := strings.ToLower(city)

	cachedWeather, err := s.cache.Get(ctx, cacheKey)
	if cachedWeather != "" && err == nil {
		var response models.WeatherResponse
		if err = json.Unmarshal([]byte(cachedWeather), &response); err == nil {
			return response, nil
		}
	}

	// Make API call
	url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?key=%s&contentType=json", city, s.apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return models.WeatherResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.WeatherResponse{}, fmt.Errorf("error: %d", resp.StatusCode)
	}

	// Decode API data into struct
	var weatherData models.WeatherData
	err = json.NewDecoder(resp.Body).Decode(&weatherData)
	if err != nil {
		return models.WeatherResponse{}, fmt.Errorf("error decoding resp body: %s", err)
	}

	// Transform weatherData struct into simpler response
	response := models.WeatherResponse{
		Location:    weatherData.Location,
		Temperature: weatherData.Days[0].Temperature,
		UVIndex:     weatherData.Days[0].UVIndex,
		Conditions:  weatherData.Days[0].Conditions,
		Time:        time.Now(),
	}

	// Cache data
	jsonData, err := json.Marshal(response)
	if err == nil {
		err = s.cache.Set(ctx, cacheKey, string(jsonData), 12*time.Hour)
		if err != nil {
			return models.WeatherResponse{}, err
		}
	}

	return response, nil
}
