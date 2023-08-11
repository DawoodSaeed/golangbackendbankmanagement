package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type APIServer struct {
	listenAddr string
	store      Storage
}

type APIError struct {
	Error string
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func makeHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			err := writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
			if err != nil {
				return
			}

		}
	}
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// To create a json web token'
func createJWT(account *Account) (string, error) {
	// Create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"number": account.Number,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	// Sign the token with a secret key
	secretKey := []byte(os.Getenv("SECRET_KEY"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println("Error creating token:", err)
		return "", err
	}

	return tokenString, nil

}

//validate the jwt token
func validateTheJWT(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("x-auth")
		fmt.Println(authHeader)
		token, tokenParseErr := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
			return []byte(""), nil
		})

		if tokenParseErr != nil {
			w.Write([]byte("Sorry, token parsing failed"))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			fmt.Println(claims["number"])
			handlerFunc(w, r)
			return
		} else {
			fmt.Println("Invalid token 😒")
			w.Write([]byte("Sorry the provided token was not valid"))
			return
		}

	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account/{id}", validateTheJWT(makeHandlerFunc(s.handleGetAccountById)))
	router.HandleFunc("/account", makeHandlerFunc(s.handleAccount))
	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		return
	}
}

//ROUTE METHODS ############################
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
		break
	case "POST":
		return s.handleCreateAccount(w, r)
		break
	case "DELETE":
		return s.handleDeleteAccount(w, r)
		break
	default:
		return fmt.Errorf("method not allowed %s", r.Method)
	}

	return nil
}

// fetch all the accounts;
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	//account := NewAccount("Dawood", "saeed")
	//writeJSON(w, http.StatusOK, account)

	accounts, err := s.store.GetAccounts()
	if err != nil {
		return writeJSON(w, http.StatusInternalServerError,
			APIError{Error: "Couldn't get the accounts"})
	}

	return writeJSON(w, http.StatusOK, accounts)
}

//create an account
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	account := new(CreateAccount)
	if err := json.NewDecoder(r.Body).Decode(account); err != nil {
		return writeJSON(w, http.StatusBadRequest, APIError{Error: "Bad request"})
	}
	newAccount := NewAccount(account.FirstName, account.LastName)

	// token
	token, tokenErr := createJWT(newAccount)
	if tokenErr != nil {
		return tokenErr
	}
	w.Header().Add("X-Auth", token)
	// insert data to the database
	err := s.store.CreateAccount(newAccount)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return writeJSON(w, http.StatusOK, newAccount)
}

//validate the token (Serves as the middleware )

//Delete an account
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

//get account by ID
func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getId(r)

	if err != nil {
		return fmt.Errorf("parameter Id wasn't converted to string")
	}

	if r.Method == "GET" {
		account, err := s.store.GetAccountById(id)
		if err != nil {
			return writeJSON(w, http.StatusInternalServerError, APIError{Error: err.Error()})
		}
		return writeJSON(w, http.StatusOK, account)

	} else if r.Method == "DELETE" {
		err := s.store.DeleteAccount(id)
		if err != nil {
			return writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
		return writeJSON(w, http.StatusOK, "Account Deleted")
	} else if r.Method == "PUT" {

		//getting the request body => account
		account := new(Account)
		jsonErr := json.NewDecoder(r.Body).Decode(account)
		if jsonErr != nil {
			return writeJSON(w, http.StatusBadRequest, jsonErr)
		}
		//Now we have both account, id update data in postgres
		err := s.store.UpdateAccount(account, id)
		//Check if the data updated successfully
		if err != nil {
			return writeJSON(w, http.StatusInternalServerError,
				APIError{Error: err.Error()})
		}
	}

	return writeJSON(w, http.StatusOK, APIError{Error: "Data Updated"})
}

//To get the parameter ID
func getId(r *http.Request) (int, error) {
	stringId := mux.Vars(r)["id"]
	id, err := strconv.Atoi(stringId)
	if err != nil {
		return 0, nil
	}
	return id, nil
}
