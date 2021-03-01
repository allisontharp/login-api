package main

import (
	"fmt"
	"testing"
)

func Test_validateCreds(t *testing.T) {
	tables := []struct {
		username string
		password string
		token    string
		isValid  bool
	}{
		{"a@a.com", "123", getToken(), false},
		{"c137@onecause.com", "#th@nH@rm#y#r!$100%D0p#", getToken(), true}, // note that you have to have this ,!
	}

	for _, table := range tables {
		login := Login{}
		login.Username = table.username
		login.Password = table.password
		login.Token = table.token
		ok, err := validateCreds(login)
		fmt.Printf("%v, %v", ok, err)
		if ok != table.isValid {
			t.Errorf("Validation Login.  Object Recieved: %v, Expected isValid: %v", table, ok)
		}
	}

}
