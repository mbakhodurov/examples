package weather_center

import (
	"errors"
	"sync"
)

var ErrCityNotFound = errors.New("city not found")

type WeatherCenter struct {
	mu      sync.RWMutex
	weahter map[string]float32
}

func NewWeatherCenter() *WeatherCenter {
	return &WeatherCenter{
		weahter: make(map[string]float32),
	}
}

func (wc *WeatherCenter) SetTemperature(city string, temp float32) {
	wc.mu.Lock()
	defer wc.mu.Unlock()

	wc.weahter[city] = temp
}

func (wc *WeatherCenter) GetTemperature(city string) (float32, error) {
	temper, ok := wc.weahter[city]
	if !ok {
		return 0, ErrCityNotFound
	}

	return temper, nil
}
