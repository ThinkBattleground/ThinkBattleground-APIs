package authentication

import (
	"encoding/json"
	"log"
	"net/http"
	"thinkbattleground-apis/config"
	"thinkbattleground-apis/constants"
	"thinkbattleground-apis/models"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// ForgotPassword godoc
// @Summary Forgot password
// @Description Send OTP for password reset
// @Tags users
// @Accept  json
// @Produce  json
// @Param email body models.Email true "User email"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/forgot-password [post]
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var email models.Email

	if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
		config.WriteResponse(w, http.StatusInternalServerError, constants.INVALID_REQUEST)
		log.Printf(constants.INVALID_REQUEST+" Error: %s\n", err)
		return
	}

	ok := ValidateEmail(w, email.Email)
	if !ok {
		return
	}

	// generate otp and set expiration time for otp
	otp := config.GenerateOTP()
	expiration := time.Now().Add(2 * time.Minute)

	var firstName, lastName string

	userData := `SELECT first_name, last_name from users WHERE email=$1`
	_ = config.DB.QueryRow(userData, email.Email).Scan(&firstName, &lastName)

	insertData := `INSERT INTO forgot_password (email, otp, otp_expires) VALUES ($1, $2, $3)`
	_, err := config.DB.Exec(insertData, email.Email, otp, expiration)

	go config.CronSchedule("forgot_password")

	if err != nil {
		config.WriteResponse(w, http.StatusInternalServerError, "Failed to update password!")
		log.Printf("Error while insert data in database: %s", err)
		return
	}

	message := "We received a request to reset your password. Please use the OTP below to proceed:"

	if err := config.SendEmail(email.Email, otp, firstName+" "+lastName, message); err != nil {
		config.WriteResponse(w, http.StatusInternalServerError, "Failed to send email")
		log.Printf("Error while send email: %s", err)

		deleteUser := `DELETE FROM forgot_password WHERE email = $1`
		_, err = config.DB.Exec(deleteUser, email.Email)
		if err != nil {
			log.Printf("Failed to delete data from forgot_password: %v\n", err)
		}
		return
	}

	resp := models.ResponseWithEmail{
		Message: constants.OTP_SENT,
		Email:   email.Email,
	}

	config.WriteResponse(w, http.StatusOK, resp)
	log.Printf(constants.OTP_SENT+"! OTP is : %s", otp)
}

// ResetPasswordAfterForgotPassword godoc
// @Summary Reset password after forgot password
// @Description Reset password after forgot password
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.ResetPassword true "Reset Password"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /user/forgot-password/reset-password [put]
func ResetPasswordAfterForgotPassword(w http.ResponseWriter, r *http.Request) {
	var data models.ResetPassword

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		config.WriteResponse(w, http.StatusBadRequest, constants.INVALID_REQUEST)
		log.Printf(constants.INVALID_REQUEST+" Error: %s\n", err)
		return
	}

	ok := ValidateEmail(w, data.Email)
	if !ok {
		return
	}

	ok = ValidatePassword(w, data.Password)
	if !ok {
		return
	}

	var verified bool
	verifiedUser := `SELECT verified FROM forgot_password WHERE email = $1`
	err := config.DB.QueryRow(verifiedUser, data.Email).Scan(&verified)

	if err != nil {
		config.WriteResponse(w, http.StatusUnauthorized, constants.USER_NOT_FOUND+" or "+constants.INVALID_OTP)
		log.Println(constants.USER_NOT_FOUND + " or " + constants.INVALID_OTP + err.Error())
		return
	}

	if !verified {
		config.WriteResponse(w, http.StatusUnauthorized, "User not verified!"+constants.USER_UNAUTHORIZED)
		log.Println("User not verified!" + constants.USER_UNAUTHORIZED)
		return
	}

	// encrypt password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	updateUser := `UPDATE users SET password=$1, verified_at=$2 WHERE email=$3`
	_ = config.DB.QueryRow(updateUser, hashedPassword, time.Now(), data.Email)

	// Delete user from temp_users
	deleteUser := `DELETE FROM forgot_password WHERE email = $1`
	_, err = config.DB.Exec(deleteUser, data.Email)
	if err != nil {
		config.WriteResponse(w, http.StatusBadRequest, "failed to delete user")
		log.Printf("Failed to delete data from forgot_password: %v\n", err)
	}

	resp := models.Response{
		Message: constants.PASSWORD_UPDATE,
	}

	config.WriteResponse(w, http.StatusOK, resp)
	log.Printf(constants.PASSWORD_UPDATE)
}
