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
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body models.RegisterRequest true "Registration details"
// @Success 200 {object} models.ResponseWithEmail
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/register [post]
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.WriteResponse(w, http.StatusBadRequest, models.Response{
			Message: constants.INVALID_REQUEST,
		})
		log.Printf(constants.INVALID_REQUEST+" Error: %s\n", err)
		return
	}

	// Check Validation For all the fields
	ok := validation(w, models.Users{
		UserName: req.UserName,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	})
	if !ok {
		return
	}

	userData := `SELECT email from users WHERE email=$1`
	rows, err := config.DB.Query(userData, req.Email)
	if err != nil {
		log.Println(constants.USER_NOT_FOUND + err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			log.Fatal(err)
		}
		if email == req.Email {
			config.WriteResponse(w, http.StatusBadRequest, models.Response{
				Message: fmt.Sprintf("Email %s is registered with other User.", email),
			})
			return
		}
	}

	// encrypt password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	// generate otp and set expiration time for otp
	otp := config.GenerateOTP()
	expiration := time.Now().Add(2 * time.Minute)

	insertUser := `INSERT INTO temp_users (user_name, email, password, role, otp, otp_expires) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = config.DB.Exec(insertUser, req.UserName, req.Email, hashedPassword, req.Role, otp, expiration)

	if err != nil {
		config.WriteResponse(w, http.StatusInternalServerError, models.Response{
			Message: "Failed to register user",
		})
		log.Printf("Error while insert data in database: %s", err)
		return
	}

	go config.CronSchedule("temp_users")

	message := "Thank you for signing up for Think Battleground. Please use the OTP below to verify your email and activate your account:"

	if err := config.SendEmail(req.Email, otp, req.UserName, message); err != nil {
		config.WriteResponse(w, http.StatusInternalServerError, models.Response{
			Message: "Failed to send email",
		})
		log.Printf("Error while send email: %s", err)

		// Delete user from temp_users
		deleteUser := `DELETE FROM temp_users WHERE email = $1`
		_, err = config.DB.Exec(deleteUser, req.Email)
		if err != nil {
			log.Printf("Failed to delete temp user: %v\n", err)
		}
		return
	}

	config.WriteResponse(w, http.StatusOK, models.ResponseWithEmail{
		Message: constants.OTP_SENT,
		Email:   req.Email,
	})
	log.Printf(constants.OTP_SENT+"! OTP is : %s", otp)
}
