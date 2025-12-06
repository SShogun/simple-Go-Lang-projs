package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CurrencyResponse struct {
	Result          string             `json:"result"`
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

func fetchRates() (*CurrencyResponse, error) {
	const apiURL = "https://v6.exchangerate-api.com/v6/c48eeafbbd058bf48a23fdb3/latest/USD"

	// 1. Make the GET request
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("Error making HTTP request: %w", err)
	}

	defer resp.Body.Close()

	// 2. Check for a successful status 200 OK code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected status code: %d", resp.StatusCode)
	}

	// 3. Read the entire response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %w", err)
	}

	var rates CurrencyResponse

	err = json.Unmarshal(body, &rates)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	// 5. Basic check to ensure the API call was successful (not just the HTTP status)
	if rates.Result != "success" {
		return nil, fmt.Errorf("API response indicated failure: %s", rates.Result)
	}

	return &rates, nil
}

func conversion(amount float64, inrconv float64) {
	fmt.Printf("The given amount in INR is: %.2f\n", amount*inrconv)

}

func main() {
	fmt.Println("Command Line Currency Converter")
	fmt.Println("Fetching lastest exchagne rates from API...")

	ratesData, err := fetchRates()

	if err != nil {
		fmt.Printf("Error fetching exchange rates: %v\n", err)
		return
	}

	fmt.Printf("Successfully loaded rates! Base currency: %s\n", ratesData.BaseCode)
	fmt.Printf("Available rates count: %d\n", len(ratesData.ConversionRates))
	fmt.Printf("Indian Rupees: %.2f\n", ratesData.ConversionRates["INR"])
	inrconv := ratesData.ConversionRates["INR"]

	var amount float64
	fmt.Print("Enter amount to be converted: ")
	_, err = fmt.Scanln(&amount)
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return
	}
	conversion(amount, inrconv)
}
