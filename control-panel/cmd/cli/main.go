package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type CreateRequest struct {
	TenantID string `json:"tenantId"`
}

// The structure we expect back FROM the API
type APIResponse struct {
	Message string `json:"message"`
	Tenant  string `json:"tenant"`
	// If you update your HTTP handler to return the password, it will map here
	Password string `json:"password,omitempty"`
}

func main() {
	// 1. Basic command routing
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "create":
		handleCreate()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func handleCreate() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Missing database name.")
		fmt.Println("Usage: dbpods create <database-name>")
		os.Exit(1)
	}

	tenantID := os.Args[2]
	apiURL := "http://localhost:8080/create-database"

	fmt.Printf(" Provisioning database [%s] on the cluster...\n", tenantID)

	// 2. Prepare the JSON payload
	reqBody := CreateRequest{TenantID: tenantID}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Failed to encode request: %v\n", err)
		os.Exit(1)
	}

	// 3. Make the HTTP POST request to your Control Plane
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("API Connection Failed: Is your Control Plane running? (%v)\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// 4. Parse and display the result
	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Server returned error code: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var apiResp APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		fmt.Printf("Failed to read server response: %v\n", err)
		os.Exit(1)
	}

	// 5. Print the success message and connection instructions
	fmt.Println("\n Success! Your database is ready.")
	fmt.Println(stringsRepeat("-", 40))
	fmt.Printf("Tenant ID : %s\n", apiResp.Tenant)

	// Print the connection string so the user can easily copy/paste it
	fmt.Println("\n Connect via TCP Proxy:")
	fmt.Printf("psql -h localhost -p 5432 -U %s -d %s_data\n", apiResp.Tenant, apiResp.Tenant)
	fmt.Println(stringsRepeat("-", 40))
}

func printUsage() {
	fmt.Println("\n dbpods CLI - Database-as-a-Service Manager")
	fmt.Println("Usage:")
	fmt.Println("  dbpods create <database-name>   Spin up a new Postgres database")
}

func stringsRepeat(char string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += char
	}
	return result
}
