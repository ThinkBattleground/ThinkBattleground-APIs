package authentication

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"thinkbattleground-apis/config"
	"thinkbattleground-apis/models"
	"unicode"
)

func validation(w http.ResponseWriter, user models.Users) bool {
	ok := true

	if len(user.FirstName) == 0 {
		config.WriteResponse(w, http.StatusBadRequest, "First Name must be not empty")
		ok = false
	}

	if len(user.LastName) == 0 {
		config.WriteResponse(w, http.StatusBadRequest, "Last Name must be not empty")
		ok = false
	}

	// checking whether genre is valid
	flag := checkStringInSlice(validRole, user.Role)
	if !flag {
		config.WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("role %s is not allowed", user.Role))
		ok = false
	}

	// validate email
	email := ValidateEmail(w, user.Email)
	if !email {
		ok = false
	}

	// validate password
	password := ValidatePassword(w, user.Password)
	if !password {
		ok = false
	}
	return ok
}

func ValidatePassword(w http.ResponseWriter, password string) bool {
	ok := true
	if len(password) == 0 {
		config.WriteResponse(w, http.StatusBadRequest, "Password must be not empty")
		ok = false
	} else {
		errors := ValidatePasswordString(password)
		if len(errors) == 0 {
			log.Println("Password is strong.")
		} else {
			config.WriteResponse(w, http.StatusBadRequest, "Password is not strong enough. Issues found:")
			fmt.Println("Password is not strong enough. Issues found:")
			for _, err := range errors {
				config.WriteResponse(w, http.StatusBadRequest, "-"+err)
				log.Println("-", err)
			}
			ok = false
		}
	}
	return ok
}

func ValidateEmail(w http.ResponseWriter, email string) bool {
	ok := true
	if len(email) == 0 {
		config.WriteResponse(w, http.StatusBadRequest, "Email must be not empty")
		ok = false
	} else {
		if !strings.Contains(email, "@") {
			config.WriteResponse(w, http.StatusBadRequest, "Please provide valid email address!")
			ok = false
		}
	}
	return ok
}

func checkStringInSlice(items []string, item string) bool {
	for _, cur := range items {
		if cur == item {
			return true
		}
	}
	return false
}

func ValidatePasswordString(password string) []string {
	var errors []string

	if len(password) < 8 {
		errors = append(errors, "Password must be at least 8 characters long.")
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		errors = append(errors, "Password must contain at least one uppercase letter.")
	}
	if !hasLower {
		errors = append(errors, "Password must contain at least one lowercase letter.")
	}
	if !hasDigit {
		errors = append(errors, "Password must contain at least one digit.")
	}
	if !hasSpecial {
		errors = append(errors, "Password must contain at least one special character.")
	}
	return errors
}
