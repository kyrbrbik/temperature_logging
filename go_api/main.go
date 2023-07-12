package main

import (
	"net/http"
	"database/sql"
	"log"
	"fmt"
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
	loadDB()
	r := gin.Default()
	r.GET("/data", getData)
	r.POST("/data", addData)
	r.Run(":8080")
}

func loadDB() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "password",
		Net:    "tcp",
		Addr:   "localhost:3306",
		DBName: "sensor",
	}

	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

func addData(c *gin.Context) {
	var newData Data

	if err := c.ShouldBindJSON(&newData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	id err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted row ID:", id)
	return nil
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
