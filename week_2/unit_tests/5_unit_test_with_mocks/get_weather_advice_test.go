package main

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/mbakhodurov/examples/week_2/unit_tests/5_unit_test_with_mocks/mocks"
	"github.com/mbakhodurov/examples/week_2/unit_tests/5_unit_test_with_mocks/weather_center"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestGetWeatherAdvice(t *testing.T) {
	type weatherCenterClientMockFunc func(t *testing.T) WeatherCenterClient

	city := gofakeit.City()

	tests := []struct {
		name                    string
		city                    string
		temperature             float32
		err                     error
		expected                string
		weatherCenterClientMock weatherCenterClientMockFunc
	}{
		{
			name:     "Температура +25 градусов",
			city:     city,
			expected: "Отличная погода для прогулок",
			err:      nil,
			weatherCenterClientMock: func(t *testing.T) WeatherCenterClient {
				mockClient := mocks.NewWeatherCenterClient(t)
				mockClient.On("GetTemperature", city).Return(float32(25), nil).Once()

				return mockClient
			},
		},
		{
			name:     "Температура -15 градусов",
			city:     city,
			expected: "Прохладно, но можно выйти на улицу",
			err:      nil,
			weatherCenterClientMock: func(t *testing.T) WeatherCenterClient {
				mockClient := mocks.NewWeatherCenterClient(t)
				mockClient.On("GetTemperature", city).Return(float32(-15), nil).Once()

				return mockClient
			},
		},
		{
			name:     "Город не найден",
			city:     city,
			expected: "",
			err:      weather_center.ErrCityNotFound,
			weatherCenterClientMock: func(t *testing.T) WeatherCenterClient {
				mockClient := mocks.NewWeatherCenterClient(t)
				mockClient.On("GetTemperature", city).Return(float32(0), weather_center.ErrCityNotFound)

				return mockClient
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := getWeatherAdvice(test.weatherCenterClientMock(t), test.city)
			require.True(t, errors.Is(err, test.err))
			require.Equal(t, test.expected, res)
		})
	}
}
