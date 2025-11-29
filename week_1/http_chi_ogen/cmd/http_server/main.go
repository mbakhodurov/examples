package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	custommiddleware "github.com/mbakhodurov/examples/week_1/http_chi_ogen/internal/middleware"
	weather_v1 "github.com/mbakhodurov/examples/week_1/http_chi_ogen/pkg/openapi/weather/v1"
)

const (
	httpPort = "8080"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type WeatherStorage struct {
	mu       sync.RWMutex
	weathers map[string]*weather_v1.Weather
}

func NewWeatherStorage() *WeatherStorage {
	return &WeatherStorage{
		weathers: make(map[string]*weather_v1.Weather),
	}
}

func (s *WeatherStorage) GetWeather(city string) *weather_v1.Weather {
	s.mu.RLock()
	defer s.mu.RUnlock()
	weather, ok := s.weathers[city]
	if !ok {
		return nil
	}
	return weather
}

func (s *WeatherStorage) UpdateWeather(city string, weather *weather_v1.Weather) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.weathers[city] = weather
}

type WeatherHandler struct {
	storage *WeatherStorage
}

func NewWeatherHandler(storage *WeatherStorage) *WeatherHandler {
	return &WeatherHandler{
		storage: storage,
	}
}

func (h *WeatherHandler) GetWeatherByCity(_ context.Context, params weather_v1.GetWeatherByCityParams) (weather_v1.GetWeatherByCityRes, error) {
	weather := h.storage.GetWeather(params.City)
	if weather == nil {
		return &weather_v1.NotFoundError{
			Code:    404,
			Message: "Weather data not found for city: " + params.City,
		}, nil
	}
	return weather, nil
}

func (h *WeatherHandler) UpdateWeatherByCity(_ context.Context, req *weather_v1.UpdateWeatherRequest, params weather_v1.UpdateWeatherByCityParams) (weather_v1.UpdateWeatherByCityRes, error) {

	weather := &weather_v1.Weather{
		City:        params.City,
		Temperature: req.Temperature,
		UpdatedAt:   time.Now(),
	}
	h.storage.UpdateWeather(params.City, weather)
	return weather, nil
}

func (h *WeatherHandler) NewError(_ context.Context, err error) *weather_v1.GenericErrorStatusCode {
	return &weather_v1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: weather_v1.GenericError{
			Code:    weather_v1.NewOptInt(http.StatusInternalServerError),
			Message: weather_v1.NewOptString(err.Error()),
		},
	}
}

func main() {
	storage := NewWeatherStorage()

	weatherHandler := NewWeatherHandler(storage)

	weatherServer, err := weather_v1.NewServer(weatherHandler)
	if err != nil {
		log.Fatalf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(custommiddleware.RequestLogger)

	// Монтируем обработчики OpenAPI
	r.Mount("/", weatherServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска сервера: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы сервера...")

	// Создаем контекст с таймаутом для остановки сервера
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке сервера: %v\n", err)
	}

	log.Println("✅ Сервер остановлен")
}
