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
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/login [post]
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		config.WriteResponse(w, http.StatusBadRequest, models.Response{
			Message: constants.INVALID_REQUEST,
		})
		log.Printf(constants.INVALID_REQUEST+" Error: %s\n", err)
		return
	}

	ok := ValidateEmail(w, loginReq.Email)
	if !ok {
		return
	}

	ok = ValidatePassword(w, loginReq.Password)
	if !ok {
		return
	}

	var user models.Users
	var password string

	userData := `SELECT user_name, password, role FROM users WHERE email=$1`
	err := config.DB.QueryRow(userData, loginReq.Email).Scan(&user.UserName, &password, &user.Role)

	if err != nil {
		config.WriteResponse(w, http.StatusBadRequest, models.Response{
			Message: fmt.Sprintf("User with  %s not exist, please register.", loginReq.Email),
		})
		log.Printf("User with  %s not exist, please register.\n", loginReq.Email)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(loginReq.Password))
	if err == nil {
		expireTime := time.Now().Add(20 * time.Minute)
		// creating and signing token for defined claims using HS256 method
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_name": user.UserName,
			"email":     loginReq.Email,
			"role":      user.Role,
			"exp":       expireTime.Unix(),
		})
		if err := config.LoadEnv(); err != nil {
			log.Println(constants.LOAD_ENV_ERROR)
			return
		}

		jwtKey := os.Getenv("JWTKEY")
		tokenString, err := token.SignedString([]byte(jwtKey))
		if err != nil {
			config.WriteResponse(w, http.StatusInternalServerError, models.Response{
				Message: "Error while setting token",
			})
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
			Message:    fmt.Sprintf("User %s logged in successfully! Token Expires within 20 Minutes", loginReq.Email),
			Token:      tokenString,
			ExpireTime: time.Now().Add(20 * time.Minute).Format(time.RFC3339),
		}

		config.WriteResponse(w, http.StatusOK, resp)
		fmt.Println("Token : ", tokenString)
	} else {
		config.WriteResponse(w, http.StatusUnauthorized, models.Response{
			Message: "Invalid Credentials",
		})
	}
}
