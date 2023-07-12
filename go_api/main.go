package main

import (
	"net/http"
	"database/sql"
	"log"
	"fmt"
	"time"
	"github.com/go-sql-driver/mysql"
	"github.com/gin-gonic/gin"
)

var db *sql.DB

type Data struct {
	ID int64
	Temperature string `json:"temperature"`
	Humidity string `json:"humidity"`
	Time string `json:"time"`
}

func main() {
	if err := loadDB(); err != nil {
		log.Fatal(err)
	}

	initTable()
	r := gin.Default()
	r.Use(logTimestamp())
	r.GET("/data", getData)
	r.POST("/data", addData)
	r.Run(":8080")
}

func loadDB() error{
	cfg := mysql.Config{
		User:   "root",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "sensor",
	}
	
	var err error
	maxRetries := 10
	retryInterval := 5 * time.Second

	for retries :=0; retries < maxRetries; retries++ {
		db, err = sql.Open("mysql", cfg.FormatDSN())
		if err != nil {
			log.Println(err)
			time.Sleep(retryInterval)
		}

		pingErr := db.Ping()
		if pingErr != nil {
			log.Println(pingErr)
			time.Sleep(retryInterval)
		}
	}
	return err
}

func addData(c *gin.Context) {
	var newData Data

	if err := c.ShouldBindJSON(&newData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	timestamp, exists := c.Get("Time")
	if !exists { 
		c.JSON(http.StatusBadRequest, gin.H{"error": "Timestamp not found"})
	}

	newData.Time = fmt.Sprintf("%v", timestamp.(int64))
	

	err := writeDB(newData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, newData)
}

func writeDB(data Data) error {
	stmt, err := db.Prepare("INSERT INTO data (temperature, humidity, time) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(data.Temperature, data.Humidity, data.Time)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
	return err
}

func getData(c *gin.Context) {
	var data Data
	rows, err := db.Query("SELECT * FROM data")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&data.ID, &data.Temperature, &data.Humidity, &data.Time)
		if err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, gin.H{"data": data})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func logTimestamp() gin.HandlerFunc {
	return func(c *gin.Context) {
		timestamp := time.Now().Unix()
		c.Set("Time", timestamp)
		c.Next()
	}
}

func initTable() {
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS data (id INTEGER PRIMARY KEY AUTO_INCREMENT, temperature FLOAT, humidity FLOAT, time INT)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res)
}
