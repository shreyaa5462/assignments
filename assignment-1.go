package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var exchangeRates = map[string]float64{
	"USD": 1.0,
	"INR": 83.12,
	"EUR": 0.93,
	"JPY": 156.82,
}

func init() {
	hour := time.Now().Hour()
	if hour < 12 {
		fmt.Println("Good Morning!")
	} else if hour < 18 {
		fmt.Println("Good Afternoon!")
	} else {
		fmt.Println("Good Evening!")
	}
}
func main() {
	args := os.Args[1:]
	if len(args) == 1 && args[0] == "--list" {
		listSupportedCurrencies()
		return
	}
	if len(args) != 3 {
		printUsage()
		os.Exit(1)
	}
	amount, sourceCurrency, targetCurrency, err := validateInput(args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		printUsage()
		os.Exit(1)
	}
	convertedAmount, err := convertCurrency(amount, sourceCurrency, targetCurrency)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%.2f %s is equivalent to %.2f %s\n", amount, sourceCurrency, convertedAmount, targetCurrency)
}
func validateInput(args []string) (float64, string, string, error) {
	// Parse amount
	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return 0, "", "", fmt.Errorf("invalid amount: '%s' is not a valid number", args[0])
	}
	if amount <= 0 {
		return 0, "", "", fmt.Errorf("amount must be a positive number")
	}
	sourceCurrency := strings.ToUpper(args[1])
	targetCurrency := strings.ToUpper(args[2])
	if _, ok := exchangeRates[sourceCurrency]; !ok {
		return 0, "", "", fmt.Errorf("unsupported source currency: %s", sourceCurrency)
	}
	if _, ok := exchangeRates[targetCurrency]; !ok {
		return 0, "", "", fmt.Errorf("unsupported target currency: %s", targetCurrency)
	}
	if sourceCurrency == targetCurrency {
		return 0, "", "", fmt.Errorf("source and target currencies cannot be the same")
	}
	return amount, sourceCurrency, targetCurrency, nil
}
func convertCurrency(amount float64, sourceCurrency, targetCurrency string) (float64, error) {
	sourceRate, ok := exchangeRates[sourceCurrency]
	if !ok {
		return 0, fmt.Errorf("internal error: source currency '%s' rate not found", sourceCurrency)
	}
	amountInUSD := 1 / sourceRate
	targetRate, ok := exchangeRates[targetCurrency]
	if !ok {
		return 0, fmt.Errorf("internal error: target currency '%s' rate not found", targetCurrency)
	}
	convertedAmount := amount * amountInUSD * targetRate
	return convertedAmount, nil
}
func listSupportedCurrencies() {
	fmt.Println("Supported Currencies and their rates (1 USD = X Currency):")
	for currency, rate := range exchangeRates {
		if currency != "USD" {
			fmt.Printf("  1 USD = %.2f %s\n", rate, currency)
		}
	}
	fmt.Println("\nTo convert 1 UNIT of other currencies to USD:")
	for currency, rate := range exchangeRates {
		if currency != "USD" {
			fmt.Printf("  1 %s = %.4f USD\n", currency, 1/rate)
		}
	}
}
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  go run main.go <amount> <source_currency> <target_currency>")
	fmt.Println("Example:")
	fmt.Println("  go run main.go 100 USD INR")
	fmt.Println("\nOptions:")
	fmt.Println("  go run main.go --list    (To list all supported currencies and their rates)")
	fmt.Println("\nSupported Currencies: USD, EUR, INR, JPY")
}

//https://github.com/shreyaa5462/assignments/pull/1
