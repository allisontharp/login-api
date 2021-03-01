package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/apex/gateway"
	"github.com/gorilla/mux"
)

// Struct for login
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// In the future, this would not be here.  Instead, you would search a database (dynamodb?) by username and return the password
func getPassword() Login {
	login := Login{}
	login.Username = "c137@onecause.com"       // this is bad
	login.Password = "#th@nH@rm#y#r!$100%D0p#" // this is worse!
	login.Token = getToken()
	return login
}

// example from asanchez.dev/blob/cors-golang-options
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// set headers
		(w).Header().Set("Access-Control-Allow-Origin", "*")
		(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		fmt.Println("ok")

		next.ServeHTTP(w, r)
		return
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(400) // bad request
		return
	}

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
// Running locally, there wasn't an issue with the timezone.
// However, after deployment, it was expecting the time to be in UTC.
// In the future, this would be handled better.
func getToken() string {
	loc, _ := time.LoadLocation("America/New_York")
	currentTime := time.Now().In(loc)
	return fmt.Sprintf("%d%d", currentTime.Hour(), currentTime.Minute())
}

func validateCreds(login Login) (bool, error) {
	// get expected username/password
	expectedLogin := getPassword()
	// these are helpful for troubleshooting
	fmt.Printf("expected: %+v\n", expectedLogin)
	fmt.Printf("actual: %+v\n", login)
	if expectedLogin == login {
		return true, nil
	} else {
		return false, errors.New("Invalid username or password")
	}
}

func main() {
	r := mux.NewRouter()
	r.Use(CORS) // handles OPTIONS which is the preflight
	r.HandleFunc("/login", loginHandler)

	http.Handle("/", r)
	log.Fatal(gateway.ListenAndServe(":10000", nil))
}
