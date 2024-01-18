package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

const secretKey = "123"

type Exhibition struct {
	ExhibitionitionID int64  `json:"exhibition_id"`
	Room              string `json:"room"`
}

func main() {
	r := gin.Default()

	r.POST("/calc_room", func(c *gin.Context) {
		var exhibition Exhibition

		if err := c.ShouldBindJSON(&exhibition); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		go func() {
			time.Sleep(3 * time.Second)
			SendStatus(exhibition)
		}()

		c.JSON(http.StatusOK, gin.H{"message": "Room generation initiated"})
	})

	r.Run(":8080")
}

func SendStatus(exhib Exhibition) bool {
	exhib.Room = generateRoom()
	url := "http://localhost:8000/api/exhibitions/" + fmt.Sprint(exhib.ExhibitionitionID) + "/update_room/"
	response, err := performPUTRequest(url, exhib)
	if err != nil {
		fmt.Println("Error sending status:", err)
		return false
	}

	if response.StatusCode == http.StatusOK {
		fmt.Println("Room sent successfully for pk:", exhib.ExhibitionitionID)
		return true
	} else {
		fmt.Println("Failed to process PUT request")
		return false
	}
}

func generateRoom() string {
	rand.Seed(time.Now().UnixNano())

	number := rand.Intn(831) + 200

	if number > 700 {
		return fmt.Sprintf("%dл", number)
	}

	if number < 400 && rand.Intn(10) < 3 {
		return fmt.Sprintf("%dэ", number)
	}

	if !strings.ContainsAny(fmt.Sprintf("%d", number), "лэ") && rand.Intn(10) < 3 {
		return fmt.Sprintf("%dл", number)
	}
	return fmt.Sprintf("%d", number)
}

func performPUTRequest(url string, data Exhibition) (*http.Response, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "Application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return resp, nil
}
