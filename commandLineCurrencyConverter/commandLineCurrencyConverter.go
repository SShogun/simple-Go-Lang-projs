package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// CurrencyResponse mirrors the JSON structure from the API.
type CurrencyResponse struct {
	Result          string             `json:"result"`
	BaseCode        string             `json:"base_code"`
	ConversionRates map[string]float64 `json:"conversion_rates"`
}

// CacheEntry holds the currency data and its fetch timestamp.
// This allows us to determine if the data is stale (expired).
type CacheEntry struct {
	Data      *CurrencyResponse
	FetchedAt time.Time
}

// Global in-memory cache variable. Initialized to nil.
var rateCache *CacheEntry = nil

// Define the maximum duration before the cache is considered expired.
const cacheExpiry = 30 * time.Minute

// fetchRates gets exchange rates, prioritizing the cache if the data is fresh.
func fetchRates() (*CurrencyResponse, error) {
	// We use the public USD base endpoint for stability in this learning project.
	const apiURL = "https://v6.exchangerate-api.com/v6/c48eeafbbd058bf48a23fdb3/latest/USD"

	// --- 1. CACHE CHECK ---
	if rateCache != nil {
		timeSinceFetch := time.Since(rateCache.FetchedAt)

		if timeSinceFetch < cacheExpiry {
			// Calculate time remaining and round it for cleaner output
			timeRemaining := cacheExpiry - timeSinceFetch
			fmt.Printf("✅ Using rates from cache (Time remaining: %v)\n", timeRemaining.Round(time.Second))
			return rateCache.Data, nil // Cache is valid, return cached data.
		}
		fmt.Println("⚠️ Cache expired. Fetching new data from API...")
	} else {
		fmt.Println("Cache is empty. Fetching initial data from API...")
	}
	// --- END CACHE CHECK ---

	// --- 2. API CALL (Only runs if cache is invalid or empty) ---
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}

	// ALWAYS close the response body to prevent resource leaks.
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var rates CurrencyResponse
	err = json.Unmarshal(body, &rates)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	if rates.Result != "success" {
		return nil, fmt.Errorf("API response indicated failure: %s", rates.Result)
	}

	// --- 3. CACHE UPDATE ---
	// Data is fresh and valid. Store it in the cache with the current time.
	rateCache = &CacheEntry{
		Data:      &rates,
		FetchedAt: time.Now(), // Record the time of the successful fetch
	}
	fmt.Println("✅ Successfully updated cache with new rates.")
	// --- END CACHE UPDATE ---

	return &rates, nil
}

// convertCurrency performs the conversion using USD as the base intermediary.
// It uses the formula: Value in Target = (Value in Source / Source Rate) * Target Rate
func convertCurrency(amount float64, from, to string, rates map[string]float64) (float64, error) {
	// Standardize inputs to uppercase, as API codes are typically uppercase.
	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	if from == to {
		return amount, nil
	}

	rateFrom, okFrom := rates[from]
	rateTo, okTo := rates[to]

	// Check if the requested currencies are available in the fetched rates.
	if !okFrom {
		return 0, fmt.Errorf("source currency not supported: %s", from)
	}
	if !okTo {
		return 0, fmt.Errorf("target currency not supported: %s", to)
	}

	// Step 1: Convert the amount FROM the source currency TO the Base currency (USD).
	// Since the rate is (1 USD = X Currency), we divide by the rate to get USD value.
	valueInUSD := amount / rateFrom

	// Step 2: Convert the value from the Base currency (USD) TO the target currency.
	convertedAmount := valueInUSD * rateTo

	return convertedAmount, nil
}

func main() {
	fmt.Println("----------------------------------")
	fmt.Println(" Command Line Currency Converter ")
	fmt.Println("----------------------------------")

	// Get rates (will use API or cache)
	ratesData, err := fetchRates()
	if err != nil {
		fmt.Printf("Fatal Error: Could not load exchange rates: %v\n", err)
		return
	}

	// --- USER INPUT ---
	var amount float64
	var fromCurrency, toCurrency string

	// Display the base currency used for conversion
	fmt.Printf("\nConversion base currency used by API: %s\n", ratesData.BaseCode)

	// Get Amount
	fmt.Print("Enter amount to be converted (e.g., 100.50): ")
	// We handle the possibility of an error from Scanln, which is good practice.
	_, err = fmt.Scanln(&amount)
	if err != nil && err.Error() != "unexpected newline" {
		fmt.Printf("Error reading amount: %v\n", err)
		return
	}

	// Get From Currency
	fmt.Print("Enter source currency code (e.g., EUR): ")
	_, err = fmt.Scanln(&fromCurrency)
	if err != nil && err.Error() != "unexpected newline" {
		fmt.Printf("Error reading source currency: %v\n", err)
		return
	}

	// Get To Currency
	fmt.Print("Enter target currency code (e.g., INR): ")
	_, err = fmt.Scanln(&toCurrency)
	if err != nil && err.Error() != "unexpected newline" {
		fmt.Printf("Error reading target currency: %v\n", err)
		return
	}

	// --- CONVERSION ---
	convertedAmount, err := convertCurrency(amount, fromCurrency, toCurrency, ratesData.ConversionRates)

	if err != nil {
		fmt.Printf("\nConversion Error: %v\n", err)
		return
	}

	// --- RESULT ---
	fmt.Println("----------------------------------")
	// Ensure result is formatted neatly
	fmt.Printf("%.2f %s = %.2f %s\n", amount, strings.ToUpper(fromCurrency), convertedAmount, strings.ToUpper(toCurrency))
	fmt.Println("----------------------------------")

	// Optional: Show the cache in action by trying to fetch again
	fmt.Println("\n--- Demonstrating Cache Use (Immediate Second Fetch) ---")
	// This second call proves the caching works because it will print the "Using rates from cache" message.
	_, _ = fetchRates()
}
