package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WeatherData struct {
	Result []struct {
		Date        string `json:"date"`
		Day         string `json:"day"`
		Icon        string `json:"icon"`
		Description string `json:"description"`
		Status      string `json:"status"`
		Degree      string `json:"degree"`
		Min         string `json:"min"`
		Max         string `json:"max"`
		Night       string `json:"night"`
		Humidity    string `json:"humidity"`
	} `json:"result"`
}

func main() {
	http.HandleFunc("/weather", getWeatherHandler)
	http.ListenAndServe(":8080", nil)
}

func getWeatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")

	if city == "" {
		http.Error(w, "City parameter is required", http.StatusBadRequest)
		return
	}

	apiKey := "1Y64fhvxv3v5rpRPhTVcpK:70MtY4bOTV8CqaYIbY0rR9"
	url := fmt.Sprintf("https://api.collectapi.com/weather/getWeather?data.lang=tr&data.city=%s", city)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, "Error creating request", http.StatusInternalServerError)
		return
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "apikey "+apiKey)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending request", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	var weatherData WeatherData
	if err := json.NewDecoder(res.Body).Decode(&weatherData); err != nil {
		http.Error(w, "Error decoding response body", http.StatusInternalServerError)
		return
	}

	responseJSON, err := json.Marshal(weatherData)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
