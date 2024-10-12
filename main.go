package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const weatherApiKey = "5b2ed480d2ac48f7b6e150214241210"
const weatherApiURL = "http://api.weatherapi.com/v1/current.json?key=" + weatherApiKey + "&q="

// Структура для данных о погоде с WeatherAPI
type WeatherData struct {
	Current struct {
		TempC float64 `json:"temp_c"` // температура в градусах Цельсия
	} `json:"current"`
}

// Структура для данных о местоположении с IP-API
type LocationData struct {
	City string `json:"city"`
}

// Функция для получения текущего города по IP
func getCity() (string, error) {
	resp, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса к IP-API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ошибка сервера IP-API: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ошибка чтения ответа IP-API: %v", err)
	}

	var locationData LocationData
	err = json.Unmarshal(body, &locationData)
	if err != nil {
		return "", fmt.Errorf("ошибка парсинга данных местоположения: %v", err)
	}

	return locationData.City, nil
}

// Функция для получения погоды по городу
func getWeather(city string) (float64, error) {
	resp, err := http.Get(weatherApiURL + city)
	if err != nil {
		return 0, fmt.Errorf("ошибка выполнения запроса к WeatherAPI: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("ошибка сервера WeatherAPI: %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("ошибка чтения ответа WeatherAPI: %v", err)
	}

	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга данных погоды: %v", err)
	}

	return weatherData.Current.TempC, nil
}

func main() {
	// Получаем город
	city, err := getCity()
	if err != nil {
		fmt.Println("Ошибка определения города:", err)
		os.Exit(1)
	}

	// Получаем погоду для этого города
	temp, err := getWeather(city)
	if err != nil {
		fmt.Println("Ошибка получения погоды:", err)
		os.Exit(1)
	}

	// Вывод погоды в системное меню
	fmt.Printf("🌡 %.1f°C in %s | color=blue\n", temp, city)

	// Добавляем разделитель и кнопку для обновления
	fmt.Println("---")
	fmt.Println("Обновить | refresh=true")
}
