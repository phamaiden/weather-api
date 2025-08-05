package main

import (
	"log"
	"os"

	"github.com/lpernett/godotenv"
	"github.com/phamaiden/weather-api/internal/cache"
	"github.com/phamaiden/weather-api/internal/handlers"
	"github.com/phamaiden/weather-api/internal/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env")
	}

	cfg := config{
		addr: os.Getenv("ADDR"),
	}

	cache := cache.NewRedisCache()

	app := &application{
		config: cfg,
		cache:  cache,
	}

	// initialize services
	ws := services.NewWeatherService(os.Getenv("APIKEY"), &cache)

	// initialize handlers
	wh := handlers.NewWeatherHandler(ws)

	// initialize mux
	r := app.mount(wh)

	// start server
	log.Fatal(app.run(r))
}
