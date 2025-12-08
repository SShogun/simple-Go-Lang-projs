package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
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

// Cache file path for persistent storage
const cacheFilePath = "currency_cache.json"

// Define the maximum duration before the cache is considered expired.
// This can be overridden via command-line flag.
var cacheExpiry = 30 * time.Minute

// List of common currency codes for validation
var validCurrencies = []string{
	"USD", "EUR", "GBP", "JPY", "AUD", "CAD", "CHF", "CNY", "INR", "MXN",
	"BRL", "ZAR", "RUB", "KRW", "SGD", "HKD", "NOK", "SEK", "DKK", "NZD",
	"TRY", "PLN", "THB", "MYR", "IDR", "PHP", "CZK", "HUF", "ILS", "AED",
}

// loadCacheFromDisk reads the cache file and populates the global cache variable.
func loadCacheFromDisk() error {
	data, err := os.ReadFile(cacheFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, not an error
		}
		return fmt.Errorf("error reading cache file: %w", err)
	}

	var entry CacheEntry
	err = json.Unmarshal(data, &entry)
	if err != nil {
		return fmt.Errorf("error parsing cache file: %w", err)
	}

	rateCache = &entry
	fmt.Println("📁 Loaded cache from disk")
	return nil
}

// saveCacheToDisk writes the current cache to a JSON file.
func saveCacheToDisk() error {
	if rateCache == nil {
		return nil // Nothing to save
	}

	data, err := json.MarshalIndent(rateCache, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling cache: %w", err)
	}

	err = os.WriteFile(cacheFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing cache file: %w", err)
	}

	return nil
}

// fetchRates gets exchange rates, prioritizing the cache if the data is fresh.
// Includes retry logic with exponential backoff.
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
		fmt.Println("💾 Cache is empty. Fetching initial data from API...")
	}
	// --- END CACHE CHECK ---

	// --- 2. API CALL WITH RETRY LOGIC (Only runs if cache is invalid or empty) ---
	var resp *http.Response
	var body []byte
	var err error
	maxRetries := 3
	baseDelay := 1 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(math.Pow(2, float64(attempt-1))) * baseDelay
			fmt.Printf("🔄 Retry attempt %d/%d after %v...\n", attempt, maxRetries, delay)
			time.Sleep(delay)
		}

		resp, err = http.Get(apiURL)
		if err != nil {
			if attempt == maxRetries {
				return nil, fmt.Errorf("error making HTTP request after %d retries: %w", maxRetries, err)
			}
			fmt.Printf("⚠️ HTTP request failed: %v\n", err)
			continue
		}

		// ALWAYS close the response body to prevent resource leaks.
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if attempt == maxRetries {
				return nil, fmt.Errorf("unexpected HTTP status code: %d", resp.StatusCode)
			}
			fmt.Printf("⚠️ Got status code %d\n", resp.StatusCode)
			continue
		}

		// Success! Read the response body
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			if attempt == maxRetries {
				return nil, fmt.Errorf("error reading response body: %w", err)
			}
			fmt.Printf("⚠️ Error reading response: %v\n", err)
			continue
		}

		break // Success, exit retry loop
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

	// Save cache to disk for persistence
	if err := saveCacheToDisk(); err != nil {
		fmt.Printf("⚠️ Warning: Could not save cache to disk: %v\n", err)
	}
	// --- END CACHE UPDATE ---

	return &rates, nil
}

// validateCurrency checks if a currency code is valid and suggests alternatives if not.
func validateCurrency(code string, rates map[string]float64) error {
	code = strings.ToUpper(code)
	if _, exists := rates[code]; exists {
		return nil
	}

	// Find similar currency codes
	var suggestions []string
	for validCode := range rates {
		if strings.HasPrefix(validCode, code[:min(1, len(code))]) {
			suggestions = append(suggestions, validCode)
			if len(suggestions) >= 5 {
				break
			}
		}
	}

	if len(suggestions) > 0 {
		return fmt.Errorf("currency code '%s' not found. Did you mean: %s?", code, strings.Join(suggestions, ", "))
	}
	return fmt.Errorf("currency code '%s' not supported", code)
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
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

	// Validate currencies with helpful error messages
	if err := validateCurrency(from, rates); err != nil {
		return 0, fmt.Errorf("source currency error: %w", err)
	}
	if err := validateCurrency(to, rates); err != nil {
		return 0, fmt.Errorf("target currency error: %w", err)
	}

	rateFrom := rates[from]
	rateTo := rates[to]

	// Step 1: Convert the amount FROM the source currency TO the Base currency (USD).
	// Since the rate is (1 USD = X Currency), we divide by the rate to get USD value.
	valueInUSD := amount / rateFrom

	// Step 2: Convert the value from the Base currency (USD) TO the target currency.
	convertedAmount := valueInUSD * rateTo

	return convertedAmount, nil
}

func main() {
	// Parse command-line flags
	cacheMinutes := flag.Int("cache", 30, "Cache expiry time in minutes")
	batchMode := flag.Bool("batch", false, "Enable batch conversion mode (multiple conversions)")
	flag.Parse()

	// Set cache expiry from flag
	cacheExpiry = time.Duration(*cacheMinutes) * time.Minute

	fmt.Println("----------------------------------")
	fmt.Println(" Command Line Currency Converter ")
	fmt.Println("----------------------------------")
	fmt.Printf("⚙️  Cache expiry: %v\n", cacheExpiry)

	// Load cache from disk if available
	if err := loadCacheFromDisk(); err != nil {
		fmt.Printf("⚠️ Warning: Could not load cache from disk: %v\n", err)
	}

	// Get rates (will use API or cache)
	ratesData, err := fetchRates()
	if err != nil {
		fmt.Printf("Fatal Error: Could not load exchange rates: %v\n", err)
		return
	}

	// Display the base currency used for conversion
	fmt.Printf("\n💱 Conversion base currency used by API: %s\n", ratesData.BaseCode)

	// Batch conversion mode
	if *batchMode {
		fmt.Println("\n🔄 Batch Mode Enabled - Enter conversions (or 'quit' to exit)")
		for {
			var amount float64
			var fromCurrency, toCurrency string

			fmt.Print("\nEnter amount (or 'quit'): ")
			var input string
			fmt.Scanln(&input)
			if strings.ToLower(input) == "quit" {
				break
			}

			_, err := fmt.Sscanf(input, "%f", &amount)
			if err != nil {
				fmt.Printf("❌ Invalid amount: %v\n", err)
				continue
			}

			fmt.Print("Enter source currency code (e.g., EUR): ")
			fmt.Scanln(&fromCurrency)

			fmt.Print("Enter target currency code (e.g., INR): ")
			fmt.Scanln(&toCurrency)

			convertedAmount, err := convertCurrency(amount, fromCurrency, toCurrency, ratesData.ConversionRates)
			if err != nil {
				fmt.Printf("\n❌ Conversion Error: %v\n", err)
				continue
			}

			fmt.Println("----------------------------------")
			fmt.Printf("✅ %.2f %s = %.2f %s\n", amount, strings.ToUpper(fromCurrency), convertedAmount, strings.ToUpper(toCurrency))
			fmt.Println("----------------------------------")
		}
		fmt.Println("\n👋 Goodbye!")
		return
	}

	// --- SINGLE CONVERSION MODE ---
	var amount float64
	var fromCurrency, toCurrency string

	// Get Amount
	fmt.Print("\nEnter amount to be converted (e.g., 100.50): ")
	_, err = fmt.Scanln(&amount)
	if err != nil {
		fmt.Printf("Error reading amount: %v\n", err)
		return
	}

	// Get From Currency
	fmt.Print("Enter source currency code (e.g., EUR): ")
	_, err = fmt.Scanln(&fromCurrency)
	if err != nil {
		fmt.Printf("Error reading source currency: %v\n", err)
		return
	}

	// Get To Currency
	fmt.Print("Enter target currency code (e.g., INR): ")
	_, err = fmt.Scanln(&toCurrency)
	if err != nil {
		fmt.Printf("Error reading target currency: %v\n", err)
		return
	}

	// --- CONVERSION ---
	convertedAmount, err := convertCurrency(amount, fromCurrency, toCurrency, ratesData.ConversionRates)

	if err != nil {
		fmt.Printf("\n❌ Conversion Error: %v\n", err)
		return
	}

	// --- RESULT ---
	fmt.Println("----------------------------------")
	fmt.Printf("✅ %.2f %s = %.2f %s\n", amount, strings.ToUpper(fromCurrency), convertedAmount, strings.ToUpper(toCurrency))
	fmt.Println("----------------------------------")

	fmt.Println("\n💡 Tip: Use -batch flag for multiple conversions or -cache=60 to set cache expiry to 60 minutes")
}
