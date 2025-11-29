package models

import "sync"

type WeatherStorage struct {
	mu sync.RWMutex

	weathers map[string]*Weather
}

func NewWeatherStorage() *WeatherStorage {
	return &WeatherStorage{
		weathers: make(map[string]*Weather),
	}
}

func (s *WeatherStorage) GetAllWeathers() []*Weather {
	s.mu.RLock()
	defer s.mu.RUnlock()

	weathers := make([]*Weather, 0, len(s.weathers))

	for _, w := range s.weathers {
		weathers = append(weathers, w)
	}
	return weathers
}

func (s *WeatherStorage) UpdateWeather(w *Weather) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.weathers[w.City] = w
}

func (s *WeatherStorage) DeleteWeather(city string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.weathers, city)
}

func (s *WeatherStorage) GetWeather(city string) *Weather {
	s.mu.RLock()
	defer s.mu.RUnlock()

	weather, ok := s.weathers[city]
	if !ok {
		return nil
	}
	return weather
}
