package main

import (
	rand2 "math/rand"
)

type Account struct {
	Id        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
}

func NewAccount(FirstName, LastName string) *Account {
	account := &Account{FirstName: FirstName, LastName: LastName, Number: int64(rand2.Intn(10000))}
	return account
}
