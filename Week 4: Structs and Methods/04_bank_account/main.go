package main

import (
	"errors"
	"fmt"
)

type BankAccount struct {
	Owner   string
	balance float64 // unexported: inaccessible from other packages
}

func (b *BankAccount) Deposit(amount float64) {
	if amount > 0 {
		b.balance += amount
	}
}

func (b *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("withdrawal amount must be positive")
	}
	if amount > b.balance {
		return errors.New("insufficient funds")
	}
	b.balance -= amount
	return nil
}

func (b *BankAccount) Balance() float64 {
	return b.balance
}

func main() {
	accounts := []BankAccount{
		{Owner: "Ali", balance: 1000},
		{Owner: "Bota", balance: 2500},
		{Owner: "Dias", balance: 500},
	}

	accounts[0].Deposit(200)
	if err := accounts[1].Withdraw(300); err != nil {
		fmt.Println("Withdrawal error:", err)
	}
	if err := accounts[2].Withdraw(1000); err != nil {
		fmt.Println("Expected withdrawal error:", err)
	}

	// Use indices, not `for _, account := range accounts`, because range
	// returns a copy of each struct. The index gives access to the real element.
	for i := range accounts {
		accounts[i].Deposit(accounts[i].Balance() * 0.05)
	}

	for _, account := range accounts {
		fmt.Printf("%s balance: %.2f\n", account.Owner, account.Balance())
	}
}
