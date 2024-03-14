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
	Status string `json:"status"`
	Data   Status `json:"data"`
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

	status := getStatus(currentData)
	response := StatusResponse{
		Status: status,
		Data:   currentData,
	}
	c.JSON(http.StatusOK, response)
}

func getStatus(data Status) string {
	if data.Water < 5 || data.Wind < 6 {
		return "Aman"
	} else if (data.Water >= 5 && data.Water <= 8) || (data.Wind >= 6 && data.Wind <= 15) {
		return "Siaga"
	} else {
		return "Bahaya"
	}
}

func updateStatus() {
	statusMutex.Lock()
	defer statusMutex.Unlock()

	currentData = Status{
		Water: rand.Intn(100) + 1,
		Wind:  rand.Intn(100) + 1,
	}

	jsonData, err := json.Marshal(currentData)
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
