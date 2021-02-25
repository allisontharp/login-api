package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Struct for login
type Login struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
	Token    string `json:"Token"`
}

// In the future, this would not be here.  Instead, you would search a database (dynamodb?) by username and return the password
func getPassword() Login {
	login := Login{}
	login.Username = "c137@onecause.com"       // this is bad
	login.Password = "#th@nH@rm#y#r!$100%D0p#" // this is worse!
	login.Token = getToken()
	return login
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	// if r.Method != http.MethodPost {
	// 	w.WriteHeader(405) // not allowed
	// 	return
	// }

	// get the body of the POST request
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading post body: %v", err)
		w.WriteHeader(500) // internal server error
		return
	}

	// try to parse the body
	var login Login
	if err = json.Unmarshal(reqBody, &login); err != nil {
		log.Printf("Error parsing body: %v", err)
		w.WriteHeader(400) // bad request
		return
	}

	// Do the login comparison
	ok, err := validateCreds(login)

	if err != nil {
		log.Printf("Login internal server error: %v", err)
		w.WriteHeader(500) // Internal server error
		return
	}

	if !ok {
		log.Printf("Invalid username or password")
		w.WriteHeader(401) // unauthorized
	}

	w.WriteHeader(200)
}

// Get the 2 digit hour (24 hour time) and 2 digit minute at time of submission
func getToken() string {
	currentTime := time.Now()
	return fmt.Sprintf("%d%d", currentTime.Hour(), currentTime.Minute())
}

func validateCreds(login Login) (bool, error) {
	// get expected username/password
	expectedLogin := getPassword()
	if expectedLogin == login {
		return true, nil
	} else {
		return false, errors.New("Invalid username or password")
	}
	// json.NewEncoder(w).Encode(login)
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}

func main() {
	http.HandleFunc("/login", loginHandler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
