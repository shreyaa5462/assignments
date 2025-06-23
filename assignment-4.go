package main

import "fmt"

type BankAccount struct {
	Owner   string
	Balance float64
}

// a value receiver is appropriate here. Changes to 'b' inside this method
// would only affect a copy, not the 'myAccount' variable in main.
func (b BankAccount) DisplayBalance() {
	fmt.Printf("%s's current balance: $%.2f\n", b.Owner, b.Balance)
}

// Deposit is a method with a POINTER RECEIVER (b *BankAccount).
// It increases the balance of the account by the given amount.
func (b *BankAccount) Deposit(amount float64) {
	if amount <= 0 {
		fmt.Printf("Deposit amount must be positive. Skipping deposit of $%.2f.\n", amount)
		return
	}
	b.Balance += amount // This modifies the original BankAccount's balance
	fmt.Printf("Deposited $%.2f. New balance: $%.2f\n", amount, b.Balance)
}

// Withdraw is a method with a POINTER RECEIVER (b *BankAccount).
// It decreases the balance, but only if sufficient funds exist.
func (b *BankAccount) Withdraw(amount float64) {
	if amount <= 0 {
		fmt.Printf("Withdrawal amount must be positive. Skipping withdrawal of $%.2f.\n", amount)
		return
	}
	if b.Balance >= amount {
		b.Balance -= amount // This modifies the original BankAccount's balance
		fmt.Printf("Withdrew $%.2f. New balance: $%.2f\n", amount, b.Balance)
	} else {
		fmt.Printf("Insufficient funds for withdrawal of $%.2f. Current balance: $%.2f.\n", amount, b.Balance)
	}
}

// TryToModifyBalance is a standalone function that demonstrates passing a BankAccount by value.
// Any changes to 'b' within this function will only affect a copy, not the original variable passed.
func TryToModifyBalance(b BankAccount, amount float64) {
	b.Balance += amount // This only modifies the *copy* of BankAccount passed in
	fmt.Printf("  (Inside TryToModifyBalance func): Balance changed to $%.2f\n", b.Balance)
}

func main() {
	fmt.Println("--- Bank Account System Demonstration ---")

	myAccount := BankAccount{
		Owner:   "Alice",
		Balance: 100.00, // Initial balance
	}
	fmt.Println("\n--- Initial State ---")
	myAccount.DisplayBalance()
	fmt.Println("\n--- Deposit Funds ---")
	myAccount.Deposit(50.50)
	myAccount.DisplayBalance() // Observe the change
	fmt.Println("\n--- Attempting Withdrawals ---")
	myAccount.Withdraw(25.00)
	myAccount.DisplayBalance()
	myAccount.Withdraw(200.00)
	myAccount.DisplayBalance()

	fmt.Println("\n--- Final State ---")
	myAccount.DisplayBalance()
	fmt.Println("\n--- Demonstrating Value Receiver vs. Pointer Receiver Effect ---")
	tempAccount := BankAccount{Owner: "Bob", Balance: 50.00}
	fmt.Printf("Bob's balance BEFORE TryToModifyBalance: $%.2f\n", tempAccount.Balance)
	TryToModifyBalance(tempAccount, 20.00) // Pass by value
	fmt.Printf("Bob's balance AFTER TryToModifyBalance: $%.2f (Did not change externally)\n", tempAccount.Balance)

	fmt.Printf("Alice's balance BEFORE Deposit (again): $%.2f\n", myAccount.Balance)
	myAccount.Deposit(10.00) // This actually modifies myAccount
	fmt.Printf("Alice's balance AFTER Deposit (again): $%.2f (Changed externally)\n", myAccount.Balance)
}
