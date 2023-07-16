package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"html/template"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Data struct {
	ID int `json:"id"`
	Temperature string `json:"temperature"`
	Humidity string `json:"humidity"`
	Time string `json:"time"`
}

func main() {

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", indexHandler)

	r.Run(":8090")
}

func fetchAPI() (string, error) {
	url := "http://temperature-api.temperature.svc.cluster.local:8080/data/last"
	//url := "http://127.0.0.1:8080/data/last"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: ", err)	
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
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
		Title string
		Result string
		Hot string
	}{
		Title: "Temperature",
		Result: result,
		Hot: hot,
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func isHot(temp string) string {

//	var result map[string]Data
//	err := json.Unmarshal([]byte(temp), &result)
//	if err != nil {
//		fmt.Println("Error: ", err)
//	}
//
//	temperature := result["data"].Temperature
//	floatTemp, err := strconv.ParseFloat(temperature, 64)
//	if err != nil {
//		fmt.Println("Error isHot: ", err)
//	}

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


