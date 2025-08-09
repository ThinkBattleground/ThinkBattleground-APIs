package authentication

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"thinkbattleground-apis/config"
	"thinkbattleground-apis/constants"
	"thinkbattleground-apis/models"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// LoginUser godoc
// @Summary Login user
// @Description Login user and return JWT token
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.Users true "User object"
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/login [post]
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user models.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		config.WriteResponse(w, http.StatusBadRequest, constants.INVALID_REQUEST)
		log.Printf(constants.INVALID_REQUEST+" Error: %s\n", err)
		return
	}

	ok := ValidateEmail(w, user.Email)
	if !ok {
		return
	}

	ok = ValidatePassword(w, user.Password)
	if !ok {
		return
	}

	var password string

	userData := `SELECT user_id, first_name, last_name, password, role FROM users WHERE email=$1`
	err := config.DB.QueryRow(userData, user.Email).Scan(&user.UserId, &user.FirstName, &user.LastName, &password, &user.Role)

	if err != nil {
		config.WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("User with  %s not exist, please register.", user.Email))
		log.Printf("User with  %s not exist, please register.\n", user.Email)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err == nil {
		expireTime := time.Now().Add(20 * time.Minute)
		// creating and signing token for defined claims using HS256 method
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":         user.UserId,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"password":   password,
			"role":       user.Role,
			"exp":        expireTime.Unix(),
		})
		if err := config.LoadEnv(); err != nil {
			log.Println(constants.LOAD_ENV_ERROR)
			return
		}

		jwtKey := os.Getenv("JWTKEY")
		tokenString, err := token.SignedString([]byte(jwtKey))
		if err != nil {
			config.WriteResponse(w, http.StatusInternalServerError, "Error while setting token")
			log.Printf("Error while setting token: %s\n", err)
			return
		}

		// setting cookie for "token"
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expireTime,
			Path:    "/",
		})

		resp := models.LoginResponse{
			Message:    fmt.Sprintf("User %s logged in successfully! Token Expires within 20 Minutes", user.Email),
			Token:      tokenString,
			ExpireTime: time.Now().Add(20 * time.Minute).Format(time.RFC3339),
		}

		config.WriteResponse(w, http.StatusOK, resp)
		fmt.Println("Token : ", tokenString)
	} else {
		config.WriteResponse(w, http.StatusUnauthorized, "Invalid Credentials")
	}
}
