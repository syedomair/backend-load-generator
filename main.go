package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {

	baseURL := getBaseURL()
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

func getBaseURL() string {
	baseURL := ""
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		log.Println("Falling back to in-cluster config")
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Error loading kube config: %v", err)
		}
	}

	// Create Kubernetes client
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating k8s client: %v", err)
	}

	// Get service information
	serviceName := "istio-ingressgateway"
	namespace := "istio-system"

	service, err := kubeClient.CoreV1().Services(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Error getting service: %v", err)
	}

	// Find the http2 port
	for _, port := range service.Spec.Ports {
		if port.Name == "http2" && port.Port == 80 {
			nodePort := port.NodePort
			nodes, err := kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Fatalf("Error listing nodes: %v", err)
			}
			if len(nodes.Items) > 0 {
				nodeIP := nodes.Items[0].Status.Addresses[0].Address
				baseURL = fmt.Sprintf("http://%s:%d\n", nodeIP, nodePort)
			}
		}
	}
	return baseURL
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
