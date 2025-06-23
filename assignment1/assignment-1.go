package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	ErrInvalidAmountInput         = errors.New("invalid amount: not a valid number")
	ErrAmountNotPositive          = errors.New("amount must be a positive number")
	ErrUnsupportedSourceCurrency  = errors.New("unsupported source currency")
	ErrUnsupportedTargetCurrency  = errors.New("unsupported target currency")
	ErrSameCurrencies             = errors.New("source and target currencies cannot be the same")
	ErrInternalSourceRateNotFound = errors.New("internal error: source currency rate not found")
	ErrInternalTargetRateNotFound = errors.New("internal error: target currency rate not found")
	ErrInvalidArgumentsCount      = errors.New("invalid number of arguments")
)

const (
	CurrencyUSD                = "USD"
	MorningGreetingThreshold   = 12
	AfternoonGreetingThreshold = 18
	ExpectedInputArgsCount     = 3
	MinimumArgsForProgram      = 2
)

type CurrencyConverter struct {
	rates map[string]float64
}

func NewCurrencyConverter() *CurrencyConverter {
	return &CurrencyConverter{
		rates: map[string]float64{
			CurrencyUSD: 1.0,
			"EUR":       0.85,
			"GBP":       0.75,
			"JPY":       110.0,
			"INR":       83.0,
		},
	}
}

func getGreeting() string {
	hour := time.Now().Hour()

	switch {
	case hour < MorningGreetingThreshold:
		return "Good Morning"
	case hour < AfternoonGreetingThreshold:
		return "Good Afternoon"
	default:
		return "Good Evening"
	}
}

func (cc *CurrencyConverter) convertCurrency(amount float64, sourceCurrency, targetCurrency string) (convertedAmount float64, err error) {
	sourceRate, ok := cc.rates[sourceCurrency]
	if !ok {
		return 0, fmt.Errorf("%w: '%s'", ErrInternalSourceRateNotFound, sourceCurrency)
	}

	targetRate, ok := cc.rates[targetCurrency]
	if !ok {
		return 0, fmt.Errorf("%w: '%s'", ErrInternalTargetRateNotFound, targetCurrency)
	}

	amountInUSD := amount / sourceRate
	convertedAmount = amountInUSD * targetRate

	return convertedAmount, nil
}

func validateInput(args []string) (amount float64, sourceCurrency, targetCurrency string, err error) {
	if len(args) != ExpectedInputArgsCount {
		return 0, "", "", fmt.Errorf("%w: expected %d, got %d", ErrInvalidArgumentsCount, ExpectedInputArgsCount, len(args))
	}

	amount, err = strconv.ParseFloat(args[0], 64)
	if err != nil {
		return 0, "", "", fmt.Errorf("%w: '%s'", ErrInvalidAmountInput, args[0])
	}

	if amount <= 0 {
		return 0, "", "", ErrAmountNotPositive
	}

	sourceCurrency = args[1]
	targetCurrency = args[2]

	converter := NewCurrencyConverter()

	if _, ok := converter.rates[sourceCurrency]; !ok {
		return 0, "", "", fmt.Errorf("%w: %s", ErrUnsupportedSourceCurrency, sourceCurrency)
	}

	if _, ok := converter.rates[targetCurrency]; !ok {
		return 0, "", "", fmt.Errorf("%w: %s", ErrUnsupportedTargetCurrency, targetCurrency)
	}

	if sourceCurrency == targetCurrency {
		return 0, "", "", ErrSameCurrencies
	}

	return amount, sourceCurrency, targetCurrency, nil
}

func (cc *CurrencyConverter) displayRates() {
	fmt.Println("Supported Currencies and their rates (1 USD = X Currency):")

	for currency, rate := range cc.rates {
		if currency != CurrencyUSD {
			fmt.Printf("1 %s = %.2f %s\n", CurrencyUSD, rate, currency)
		}
	}

	fmt.Println("\nTo convert 1 UNIT of other currencies to USD:")

	for currency, rate := range cc.rates {
		if currency != CurrencyUSD {
			fmt.Printf("1 %s = %.2f %s\n", currency, 1/rate, CurrencyUSD)
		}
	}
}

func main() {
	fmt.Printf("%s! Welcome to the Currency Converter.\n", getGreeting())

	converter := NewCurrencyConverter()

	if len(os.Args) < MinimumArgsForProgram {
		fmt.Println("Usage: go run assignment-1.go <amount> <source_currency> <target_currency>")

		converter.displayRates()

		return
	}

	amount, sourceCurrency, targetCurrency, err := validateInput(os.Args[1:])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	convertedAmount, err := converter.convertCurrency(amount, sourceCurrency, targetCurrency)
	if err != nil {
		fmt.Printf("Error converting currency: %v\n", err)
		return
	}

	fmt.Printf("%.2f %s is equivalent to %.2f %s\n", amount, sourceCurrency, convertedAmount, targetCurrency)
}
