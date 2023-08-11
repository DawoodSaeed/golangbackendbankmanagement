package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(account *Account, id int) error
	GetAccountById(int) (*Account, error)
	DeleteAccount(int) error
	GetAccounts() ([]*Account, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStorage, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &PostgresStorage{
		db: db,
	}, nil
}

func (s *PostgresStorage) Init() {
	err := s.CreateAccountTable()
	if err != nil {
		log.Fatal("Couldn't create a account table")
		return
	}
	log.Fatal("Table was created successfully")
}

func (s *PostgresStorage) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS accounts_table (
		id SERIAL PRIMARY KEY, 
		first_name VARCHAR (50),
		last_name VARCHAR (50),
    	Number serial,
    	balance serial,
    	created_at timestamp
	)`

	_, err := s.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// Implementation of the interface Methods;

func (s *PostgresStorage) CreateAccount(account *Account) error {
	query := `INSERT INTO accounts_table 
    (
     first_name,
     last_name,
     number,
     balance,
     created_at
     ) VALUES ($1, $2, $3, $4, $5)`

	response, err := s.db.Exec(query,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := response.RowsAffected()
	if err != nil {
		return err
	}

	fmt.Printf("New account created with ID: %d\n", rowsAffected)
	return nil
}

func (s *PostgresStorage) UpdateAccount(account *Account, id int) error {
	query := `UPDATE accounts_table 
			  SET first_name = $1, last_name = $2
			  WHERE id = $3`

	_, err := s.db.Query(query, account.FirstName, account.LastName, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) GetAccountById(id int) (*Account, error) {
	rows, err := s.db.Query(`SELECT * FROM accounts_table where id = $1`, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return rowScanner(rows)
	}

	return nil, fmt.Errorf("sorry the account with the id of %d was not found", id)
}

func (s *PostgresStorage) DeleteAccount(id int) error {
	_, err := s.db.Query(`DELETE FROM accounts_table where id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM accounts_table`
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	var accounts []*Account
	for rows.Next() {
		account, err := rowScanner(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	fmt.Println(accounts)
	return accounts, nil
}

func rowScanner(rows *sql.Rows) (*Account, error) {
	account := new(Account)

	err := rows.Scan(&account.Id,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return account, nil

}
