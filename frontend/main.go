package main

import (
	"encoding/json"
	"fmt"
	"html/template"
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
		fmt.Println("No .env file found")
	}

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.Static("/static", "./static")

	r.GET("/", indexHandler)

	r.GET("/refresh", refreshTemperature)

	r.Run(":8090")
}

func fetchAPI() (string, error) {

	url := os.Getenv("API_URL")

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Println("Response: ", string(body))
	return string(body), err
}

func html(w http.ResponseWriter, result float64) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, result)
}

func indexHandler(c *gin.Context) {
	result := getTemp()
	hot := isHot(result)
	data := struct {
		Title     string
		Result    string
		Hot       string
		ImagePath string
	}{
		Title:     "Je u Honzy vedro?",
		Result:    result + "°C",
		Hot:       hot,
		ImagePath: "/static/images/thisisfine.jpg",
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func isHot(temp string) string {

	floatTemp, err := strconv.ParseFloat(temp, 64)
	if err != nil {
		fmt.Println("Error isHot: ", err)
	}
	fmt.Println("Temperature: ", floatTemp)
	if floatTemp <= 29.9 {
		return "Not as hot"
	} else {
		return "Hot as fuck"
	}
}

func getTemp() string {

	data, err := fetchAPI()
	if err != nil {
		fmt.Println("Error: ", err)
	}
	var result map[string]Data
	err2 := json.Unmarshal([]byte(data), &result)
	if err2 != nil {
		fmt.Println("Error getTemp: ", err)
	}
	temperature := result["data"].Temperature

	return temperature
}

func refreshTemperature(c *gin.Context) {

	result := getTemp()
	time.Sleep(2 * time.Second) //for dramatic effect
	c.String(http.StatusOK, result+"°C")
}
