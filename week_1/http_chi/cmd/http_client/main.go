package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mbakhodurov/examples/week_1/http_chi/pkg/models"
)

const (
	serverURL         = "http://localhost:8080"
	weatherAPIPath    = "/api/v1/weather/%s"
	weatherAPIAllPath = "/api/v1/weather"
	contentTypeHeader = "Content-Type"
	contentTypeJSON   = "application/json"
	requestTimeout    = 5 * time.Second
	defaultCityName   = "Moscow"
	defaultMinTemp    = -10
	defaultMaxTemp    = 40
)

var httpClient = &http.Client{
	Timeout: requestTimeout,
}

func deleteWeather(ctx context.Context, city string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, fmt.Sprintf("%s%s/%s", serverURL, weatherAPIPath, city), nil)
	if err != nil {
		return fmt.Errorf("создание DELETE-запроса: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("выполнение DELETE-запроса: %w", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("ошибка закрытия тела ответа: %v\n", cerr)
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("запись с ID %s не найдена", city)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("удаление записи (статус %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func updateWeather(ctx context.Context, city string, weather *models.Weather) (*models.Weather, error) {
	// Кодируем данные о погоде в JSON
	jsonData, err := json.Marshal(weather)
	if err != nil {
		return nil, fmt.Errorf("кодирование JSON: %w", err)
	}

	// Создаем PUT-запрос с контекстом
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		fmt.Sprintf("%s"+weatherAPIPath, serverURL, city),
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("создание PUT-запроса: %w", err)
	}
	req.Header.Set(contentTypeHeader, contentTypeJSON)

	// Выполняем запрос
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("выполнение PUT-запроса: %w", err)
	}
	defer func() {
		cerr := resp.Body.Close()
		if cerr != nil {
			log.Printf("ошибка закрытия тела ответа: %v\n", cerr)
			return
		}
	}()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("чтение тела ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("обновление данных о погоде (статус %d): %s", resp.StatusCode, string(body))
	}

	// Декодируем ответ
	var updatedWeather models.Weather
	err = json.Unmarshal(body, &updatedWeather)
	if err != nil {
		return nil, fmt.Errorf("декодирование JSON: %w", err)
	}

	return &updatedWeather, nil
}

func getAllWeather(ctx context.Context) ([]*models.Weather, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, serverURL+weatherAPIAllPath, nil)

	if err != nil {
		return nil, fmt.Errorf("создание GET-запроса: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("выполнение GET-запроса: %w", err)
	}
	defer func() {
		cerr := resp.Body.Close()
		if cerr != nil {
			log.Printf("ошибка закрытия тела ответа: %v\n", cerr)
			return
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	// Читаем тело ответа для любого ответа (для логирования при ошибке)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("чтение тела ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("получение данных о погоде (статус %d): %s", resp.StatusCode, string(body))
	}

	var weather []*models.Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, fmt.Errorf("декодирование JSON: %w", err)
	}

	return weather, nil
}

func getWeather(ctx context.Context, city string) (*models.Weather, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s"+weatherAPIPath, serverURL, city), nil)
	if err != nil {
		return nil, fmt.Errorf("создание запроса: %w", err)
	}
	resp, err := httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("выполнение GET-запроса: %w", err)
	}

	defer func() {
		cerr := resp.Body.Close()
		if cerr != nil {
			fmt.Println("Ошибка закрытия тела ответа:", cerr)
			return
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("чтение тела ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("получение данных о погоде (статус %d): %s", resp.StatusCode, string(body))
	}

	var weather models.Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		return nil, fmt.Errorf("декодирование JSON: %w", err)
	}

	return &weather, nil
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	ctx := context.Background()
	for {
		fmt.Println("\n=== Weather Menu ===")
		fmt.Println("1. Получить погоду по городу")
		fmt.Println("2. Получить все города")
		fmt.Println("3. Создать запись о погоде")
		fmt.Println("4. Обновить запись о погоде")
		fmt.Println("5. Удалить запись о погоде")
		fmt.Println("0. Выход")
		fmt.Print("Выберите пункт: ")

		if !reader.Scan() {
			break
		}
		choice := strings.TrimSpace(reader.Text())

		switch choice {
		case "1":
			fmt.Print("Введите название города: ")
			reader.Scan()
			city := strings.TrimSpace(reader.Text())

			weather, err := getWeather(ctx, city)
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else if weather == nil {
				fmt.Println("Данных нет")
			} else {
				fmt.Printf("Погода в %s: %+v\n", city, weather)
			}

		case "2":
			weathers, err := getAllWeather(ctx)
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				for _, w := range weathers {
					fmt.Print("Город: ", w.City, ", Температура: ", w.Temperature, "\n")
				}
			}

		case "3":
			fmt.Print("Введите название города: ")
			// reader.
			reader.Scan()
			city := strings.TrimSpace(reader.Text())

			fmt.Print("Введите температуру: ")
			reader.Scan()
			tempStr := strings.TrimSpace(reader.Text())
			temp, _ := strconv.ParseFloat(tempStr, 64)

			weather := &models.Weather{
				City:        city,
				Temperature: temp,
			}

			created, err := updateWeather(ctx, city, weather) // или createWeather если есть
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Создано:", created)
			}

		case "4":
			fmt.Print("Введите название города: ")
			reader.Scan()
			city := strings.TrimSpace(reader.Text())

			fmt.Print("Введите новую температуру: ")
			reader.Scan()
			tempStr := strings.TrimSpace(reader.Text())
			temp, _ := strconv.ParseFloat(tempStr, 64)

			weather := &models.Weather{
				City:        city,
				Temperature: temp,
			}

			updated, err := updateWeather(ctx, city, weather)
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Обновлено:", updated)
			}

		case "5":
			fmt.Print("Введите название города: ")
			reader.Scan()
			city := strings.TrimSpace(reader.Text())

			err := deleteWeather(ctx, city)
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Println("Удалено успешно")
			}

		case "0":
			fmt.Println("Выход...")
			return

		default:
			fmt.Println("Неверный пункт меню")
		}
	}
}
