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

func frontend() {
	
	data, err := fetchAPI()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	// parse value from json
	var result float64
	json.Unmarshal([]byte(data), &result)
	fmt.Println("Result: ", result)

	if result > 30 {
		fmt.Println("Hot")
	} else {
		fmt.Println("Cold")
	}
}

func html(w http.ResponseWriter, result float64) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, result)
}

func indexHandler(c *gin.Context) {
	result, err := fetchAPI()
	hot := isHot(result)
	if err != nil {
		fmt.Println("Error: ", err)
	}
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
	type Temperature struct {
		Temp string `json:"temp"`
	}
	err := json.Unmarshal([]byte(temp), &temp)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	unmarshalTemp := Temperature{temp}
	floatTemp, err := strconv.ParseFloat(unmarshalTemp.Temp, 64)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	if floatTemp <= 29.9 {
		return "Not as hot"
	} else {
		return "Hot as fuck"
	}
}
