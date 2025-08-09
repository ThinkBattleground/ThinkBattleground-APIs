package authentication

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"thinkbattleground-apis/config"
	"thinkbattleground-apis/constants"
	"thinkbattleground-apis/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	validRole = []string{"admin", "faculty", "student"} // valid roles
)

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user in the system
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.Users true "User object"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /users/register [post]
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.Users

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		config.WriteResponse(w, http.StatusInternalServerError, constants.INVALID_REQUEST)
		log.Printf(constants.INVALID_REQUEST+" Error: %s\n", err)
		return
	}

	// Check Validation For all the fields
	ok := validation(w, user)
	if !ok {
		return
	}

	userData := `SELECT email from users WHERE email=$1`
	rows, err := config.DB.Query(userData, user.Email)
	if err != nil {
		log.Println(constants.USER_NOT_FOUND + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Fatal(err)
		}
		if email == user.Email {
			config.WriteResponse(w, http.StatusBadRequest, fmt.Sprintf("Email %s is registered with other User.\n", email))
			return
		}
	}

	// encrypt password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// generate otp and set expiration time for otp
	otp := config.GenerateOTP()
	expiration := time.Now().Add(2 * time.Minute)

	insertUser := `INSERT INTO temp_users (first_name, last_name, email, password, role, otp, otp_expires) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = config.DB.Exec(insertUser, user.FirstName, user.LastName, user.Email, hashedPassword, user.Role, otp, expiration)

	if err != nil {
		config.WriteResponse(w, http.StatusInternalServerError, "Failed to register user")
		log.Printf("Error while insert data in database: %s", err)
		return
	}

	go config.CronSchedule("temp_users")

	message := "Thank you for signing up for Think Battleground. Please use the OTP below to verify your email and activate your account:"

	if err := config.SendEmail(user.Email, otp, user.FirstName+" "+user.LastName, message); err != nil {
		config.WriteResponse(w, http.StatusInternalServerError, "Failed to send email")
		log.Printf("Error while send email: %s", err)

		// Delete user from temp_users
		deleteUser := `DELETE FROM temp_users WHERE email = $1`
		_, err = config.DB.Exec(deleteUser, user.Email)
		if err != nil {
			log.Printf("Failed to delete temp user: %v\n", err)
		}
		return
	}

	resp := models.ResponseWithEmail{
		Message: constants.OTP_SENT,
		Email:   user.Email,
	}

	config.WriteResponse(w, http.StatusOK, resp)
	log.Printf(constants.OTP_SENT+"! OTP is : %s", otp)
}
