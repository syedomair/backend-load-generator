package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
)

func main() {
	baseURL := getEnv("BASE_URL", "http://localhost:8080")
	interval := getEnvDuration("INTERVAL_SECONDS", 3)

	client := resty.New()

	for {
		// Call /users
		resp1, err1 := client.R().Get(baseURL + "/users")
		if err1 != nil {
			log.Println("Error calling /users:", err1)
		} else {
			fmt.Printf("GET /users [%d]\n", resp1.StatusCode())
		}

		// Call /departments
		resp2, err2 := client.R().Get(baseURL + "/departments")
		if err2 != nil {
			log.Println("Error calling /departments:", err2)
		} else {
			fmt.Printf("GET /departments [%d]\n", resp2.StatusCode())
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
