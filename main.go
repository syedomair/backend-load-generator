package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

func main() {
	baseURL := getEnv("BASE_URL", "http://192.168.49.2:32273")
	interval := getEnvDuration("INTERVAL_SECONDS", 3)

	client := resty.New()

	for {
		headers := map[string]string{}

		if rand.Intn(2) == 0 {
			headers["x-user-type"] = "beta"
			fmt.Println("Sending with header: x-user-type: beta")
		} else {
			fmt.Println("Sending without header")
		}

		resp1, err1 := client.R().
			SetHeaders(headers).
			Get(baseURL + "/user/users")
		if err1 != nil {
			log.Println("Error calling /users:", err1)
		} else {
			var bodyInterface map[string]interface{}
			err := json.Unmarshal(resp1.Body(), &bodyInterface)
			if err != nil {
				log.Printf("Error unmarshaling JSON: %v\n", err)
				return
			}

			if result, ok := bodyInterface["result"]; ok {
				jsonResult, err := json.Marshal(result)
				if err != nil {
					log.Printf("Error marshaling result JSON: %v\n", err)
					return
				}
				fmt.Println(string(jsonResult))
			}

			if data, ok := bodyInterface["data"]; ok {
				jsonData, err := json.Marshal(data)
				if err != nil {
					log.Printf("Error marshaling data JSON: %v\n", err)
					return
				}
				fmt.Println(string(jsonData))
			}
		}

		// Call /departments
		resp2, err2 := client.R().Get(baseURL + "/department/departments")
		if err2 != nil {
			log.Println("Error calling /departments:", err2)
		} else {
			fmt.Printf("GET v1/departments [%d]\n", resp2.StatusCode())
		}

		time.Sleep(interval)
	}
}

func getEnv(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvDuration(key string, fallbackSeconds int) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return time.Duration(fallbackSeconds) * time.Second
	}
	d, err := time.ParseDuration(val + "s")
	if err != nil {
		return time.Duration(fallbackSeconds) * time.Second
	}
	return d
}
