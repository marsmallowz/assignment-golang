package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

type Status struct {
	Water int `json:"water"`
	Wind  int `json:"wind"`
}

type StatusResponse struct {
	WaterStatus string `json:"waterStatus"`
	WindStatus  string `json:"windStatus"`
	Data        Status `json:"data"`
}

var (
	statusMutex sync.RWMutex
	currentData Status
)

func main() {
	// Inisialisasi status awal
	updateStatus()

	router := gin.Default()

	router.GET("/status", getStatusHandler)
	router.StaticFile("/", "./index.html")

	// Inisialisasi cron job untuk memperbarui status setiap 15 detik
	c := cron.New()
	_ = c.AddFunc("@every 15s", updateStatus)
	c.Start()

	router.Run(":8080")
}

func getStatusHandler(c *gin.Context) {
	statusMutex.RLock()
	defer statusMutex.RUnlock()

	response := getStatusResponse(currentData)

	c.JSON(http.StatusOK, response)
}

func getStatusResponse(data Status) StatusResponse {
	windStatus, waterStatus := getStatus(currentData)
	response := StatusResponse{
		WindStatus:  windStatus,
		WaterStatus: waterStatus,
		Data:        currentData,
	}
	return response
}

func getStatus(data Status) (windStatus string, waterStatus string) {
	if data.Water <= 5 {
		waterStatus = "Aman"
	} else if data.Water >= 6 && data.Water <= 8 {
		waterStatus = "Siaga"
	} else {
		waterStatus = "Bahaya"
	}

	if data.Wind <= 6 {
		windStatus = "Aman"
	} else if data.Wind >= 7 && data.Wind <= 15 {
		windStatus = "Siaga"
	} else {
		windStatus = "Bahaya"
	}

	return windStatus, waterStatus
}

func updateStatus() {
	statusMutex.Lock()
	defer statusMutex.Unlock()

	currentData = Status{
		Water: rand.Intn(100) + 1,
		Wind:  rand.Intn(100) + 1,
	}

	response := getStatusResponse(currentData)
	jsonData, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonData))
	saveStatusToFile(currentData)
}

func saveStatusToFile(data Status) {
	file, err := os.Create("status.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		fmt.Println("Error encoding data:", err)
	}
}
