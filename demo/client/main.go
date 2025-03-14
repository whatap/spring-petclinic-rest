package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

type HealthCheckResponse struct {
	Status string `json:"status"`
}

type Owner struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Telephone string `json:"telephone"`
	ID        int    `json:"id"`
}

type OwnerInput struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Address   string `json:"address"`
	City      string `json:"city"`
	Telephone string `json:"telephone"`
}

type HTTPClient struct {
	BaseURL string
	Client  *http.Client
}

func NewHTTPClient(baseURL string) *HTTPClient {
	return &HTTPClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (hc *HTTPClient) HealthCheck() bool {
	url := fmt.Sprintf("%s/petclinic/actuator/health", hc.BaseURL)
	failures := 0

	for {
		resp, err := hc.Client.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			var health HealthCheckResponse
			if err := json.NewDecoder(resp.Body).Decode(&health); err == nil && health.Status == "UP" {
				fmt.Println("Health check passed. Starting queries...")
				return true
			}
		}

		failures++
		fmt.Println("Health check failed")
		if failures >= 24 {
			fmt.Println("Health check failed after 2 minutes. Exiting.")
			os.Exit(1)
		}

		time.Sleep(5 * time.Second)
	}
}

func (hc *HTTPClient) GetOwners() []Owner {
	url := fmt.Sprintf("%s/petclinic/api/owners", hc.BaseURL)

	for {
		resp, err := hc.Client.Get(url)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			time.Sleep(5 * time.Second)
			continue
		}

		var owners []Owner
		if err := json.NewDecoder(resp.Body).Decode(&owners); err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		fmt.Printf("Fetched %d owners.\n", len(owners))
		time.Sleep(5 * time.Second)
		return owners
	}
}

func (hc *HTTPClient) EditOwner(owner Owner) {
	url := fmt.Sprintf("%s/petclinic/api/owners/%d", hc.BaseURL, owner.ID)
	input := OwnerInput{
		FirstName: owner.FirstName,
		LastName:  owner.LastName,
		Address:   owner.Address,
		City:      owner.City,
		Telephone: owner.Telephone,
	}

	hc.updateOwner(http.MethodPut, url, input)
}

func (hc *HTTPClient) updateOwner(method, url string, owner OwnerInput) {
	body, err := json.Marshal(owner)
	if err != nil {
		fmt.Printf("Error marshaling request body: %s\n", err.Error())
		return
	}
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("Error creating request: %s\n", err.Error())
		return
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := hc.Client.Do(req)
	if err != nil {
		fmt.Printf("Request failed: %s\n", err.Error())
		return
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %s\n", err.Error())
		return
	}

	if resp.StatusCode < 300 {
		// fmt.Printf("Successfully updated owner ID %d\n", owner.ID)
		return
	}

	fmt.Printf("Update failed: [%s] %s\n", resp.Status, string(responseBody))
}

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var rng = rand.New(rand.NewSource(time.Now().UnixNano())) // 로컬 난수 생성기
// RandomString generates a random string of given length n
func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func (hc *HTTPClient) CreateOwner() {
	url := fmt.Sprintf("%s/petclinic/api/owners", hc.BaseURL)
	input := OwnerInput{
		FirstName: RandomString(100),
		LastName:  RandomString(100),
		Address:   RandomString(100),
		City:      RandomString(100),
		Telephone: RandomString(100),
	}

	for i := 0; i < 100; i++ {
		hc.updateOwner(http.MethodPost, url, input)
	}
}

func main() {
	baseURL := os.Getenv("PETCLINIC_URL")
	if baseURL == "" {
		fmt.Println("Environment variable PETCLINIC_URL is not set.")
		os.Exit(1)
	}

	client := NewHTTPClient(baseURL)

	// Perform health check before proceeding
	if !client.HealthCheck() {
		os.Exit(1)
	}

	numWorkersStr := os.Getenv("NUM_WORKERS")
	numWorkers, err := strconv.Atoi(numWorkersStr)
	if err != nil || numWorkers <= 0 {
		numWorkers = 5 // Default to 5 workers if not set or invalid
	}

	itemCountStr := os.Getenv("ITEM_COUNT")
	itemCount, err := strconv.Atoi(itemCountStr)
	if err != nil || itemCount <= 0 {
		numWorkers = 1000000 // Default to 5 workers if not set or invalid
	}

	var wg sync.WaitGroup

	// Start GET requests continuously in workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				owners := client.GetOwners()
				if len(owners) < itemCount {
					client.CreateOwner()
					continue
				}
				for _, owner := range owners {
					client.EditOwner(owner)
				}
			}
		}()
	}

	wg.Wait()
}
