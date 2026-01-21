package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	start := time.Now()
	godotenv.Load()

	cities := []string{"coccaglio", "riva del garda", "roma", "napoli", "venezia"}

	ch := make(chan string)
	var wg sync.WaitGroup

	for _, city := range cities {
		wg.Add(1)
		go fetchWeather(city, ch, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for result := range ch {
		fmt.Println(result)
	}

	fmt.Println("Execution time:", time.Since(start))

}

type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Temp float64 `json:"temp"`
	} `json:"main"`
}

func fetchWeather(city string, ch chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	var data WeatherData
	API_KEY := os.Getenv("OPEN_WEATHER_API_KEY")
	url := "http://api.openweathermap.org/data/2.5/weather?q=" + url.QueryEscape(city) + "&appid=" + API_KEY + "&units=metric"
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Error decoding data %s: %s\n", city, err)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		ch <- fmt.Sprintf("Error decoding data %s: %s\n", city, err)

	}

	ch <- fmt.Sprintf("City: %s, Temp: %.2f°C", data.Name, data.Main.Temp)

}
