package main

import (
	rand2 "math/rand"
	"time"
)

// CreateAccount struct => This is what comes from the user as JSON
type CreateAccount struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName""`
}

// Account struct => This is what I send as a JSON
type Account struct {
	Id        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(FirstName, LastName string) *Account {
	account := &Account{
		FirstName: FirstName,
		LastName:  LastName,
		Number:    int64(rand2.Intn(10000)),
		CreatedAt: time.Now().UTC(),
	}
	return account
}
