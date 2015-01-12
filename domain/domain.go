package domain

import (
	"time"
	"errors"
)

type DebitorRepository interface {
	FindById(id int) Debitor
}

type Debitor struct {
	ID       int
	Name     string
	Accounts []Account
}

type AccountRepository interface {
	FindById(id int) Account
}

type Account struct {
	ID       int
	Name     string
	Payments []Payment
}

type PaymentRepository interface {
	FindById(id int) Payment
}

type Payment struct {
	ID      int
	Amount  float64
	Balance float64
	Date    time.Time
}

func (debitor *Debitor) Add(account Account) error {
	debitor.Accounts = append(debitor.Accounts, account)
	return nil
}

func (account *Account) Add(payment Payment) error {
	account.Payments = append(account.Payments, payment)
	return nil
}

func (account *Account) CurrentBalance() float64 {
	lastPayment, err := account.LastPayment()
	if err != nil {
		return 0
	}
	return lastPayment.Balance
}

func (account *Account) LastPayment() (Payment, error) {
	if len(account.Payments) == 0 {
		return Payment{}, errors.New("Account has no payments")
	}

	lastPayment := account.Payments[0]
	for _, payment := range account.Payments {
		if payment.Date.After(lastPayment.Date) {
			lastPayment = payment
		}
	}

	return lastPayment, nil
}
