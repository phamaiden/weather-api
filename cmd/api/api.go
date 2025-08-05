package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/phamaiden/weather-api/internal/cache"
	"github.com/phamaiden/weather-api/internal/handlers"
)

type application struct {
	config config
	cache  cache.Cache
}

type config struct {
	addr string
}

func (app *application) mount(wh *handlers.WeatherHandler) http.Handler {
	mux := mux.NewRouter()

	mux.Use(jsonContentTypeMiddleware)
	mux.Use(limiterMiddleware)
	mux.HandleFunc("/v1/weather", wh.GetWeatherByCity).Methods("GET")
	mux.HandleFunc("/v1/weather/", wh.GetWeatherByCity).Methods("GET")
	mux.HandleFunc("/v1/weather/{city}", wh.GetWeatherByCity).Methods("GET")

	return mux
}

func (app *application) run(r http.Handler) error {

	srv := &http.Server{
		Addr:    app.config.addr,
		Handler: r,
	}

	log.Printf("server is running on %s", app.config.addr)

	return srv.ListenAndServe()
}
