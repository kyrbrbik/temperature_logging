package main

import (
	"encoding/json"
	"html/template"
	"log"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

type Data struct {
	ID          int    `json:"id"`
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
	Time        string `json:"time"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/static", "./static")

	r.GET("/", indexHandler)

	r.GET("/refresh", refreshTemperature)

	r.Run(":8090")
}

func fetchAPI(path string) (string, error) {

	url := os.Getenv("API_URL") + path

	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error: ", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error: ", err)
	}

	log.Println("Response: ", string(body))
	return string(body), err
}

func html(w http.ResponseWriter, result float64) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, result)
}

func indexHandler(c *gin.Context) {
	result := getTemp()
	hot := isHot(result)
	currentTemperature := getCurrentTemperature()
	data := struct {
		Title     string
		Result    string
		Hot       string
		ImagePath string
		CurrTemp  string
	}{
		Title:     "Je u Honzy vedro?",
		Result:    result + "°C",
		Hot:       hot,
		ImagePath: "/static/images/thisisfine.jpg",
		CurrTemp:  currentTemperature + "°C",
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func isHot(temp string) string {

	floatTemp, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		log.Println("Error isHot: ", err)
	}
	log.Println("Temperature: ", floatTemp)
	if floatTemp <= 29.9 {
		return "Not as hot"
	} else {
		return "Hot as fuck"
	}
}

func getTemp() string {

	data, err := fetchAPI("/data")
	if err != nil {
		log.Println("Error: ", err)
	}
	var result map[string]Data
	err2 := json.Unmarshal([]byte(data), &result)
	if err2 != nil {
		log.Println("Error getTemp: ", err)
	}
	temperature := result["data"].Temperature

	return temperature
}

func refreshTemperature(c *gin.Context) {

	result := getTemp()
	time.Sleep(2 * time.Second) //for dramatic effect
	c.String(http.StatusOK, result+"°C")
}

func getCurrentTemperature() string {

	data, err := fetchAPI("/temperature")
	if err != nil {
		log.Println("Error: ", err)
	}

	var result map[string]string

	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		log.Println(err)
	}

	temperature, ok := result["temperature"]
	if !ok {
		log.Println("Key temperature not found")
	}
	
	return temperature
}
