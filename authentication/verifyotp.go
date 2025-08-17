package authentication

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"thinkbattleground-apis/config"
	"thinkbattleground-apis/constants"
	"thinkbattleground-apis/models"
	"time"
)

// VerifyOTPHandler godoc
// @Summary Verify OTP for registration
// @Description Verify OTP and finalize user registration
// @Tags Users
// @Accept  json
// @Produce  json
// @Param verifyRequest body models.VerifyRequest true "Verify OTP request"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/verify-otp [post]
func VerifyOTPHandler(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.WriteResponse(w, http.StatusBadRequest, models.Response{
			Message: constants.INVALID_REQUEST,
		})
		log.Println(constants.INVALID_REQUEST + err.Error())
		return
	}

	ok := ValidateEmail(w, req.Email)
	if !ok {
		return
	}

	if len(req.OTP) != 6 {
		config.WriteResponse(w, http.StatusBadRequest, models.Response{
			Message: constants.INVALID_OTP,
		})
		log.Println(constants.INVALID_OTP)
		return
	}

	// store user details in variables
	var userName, password, role, storedOTP string
	var otpExpiry time.Time

	userData := `SELECT user_name, password, role, otp, otp_expires FROM temp_users WHERE email = $1`
	err := config.DB.QueryRow(userData, req.Email).Scan(&userName, &password, &role, &storedOTP, &otpExpiry)

	if err != nil {
		config.WriteResponse(w, http.StatusUnauthorized, models.Response{
			Message: constants.USER_NOT_FOUND + " or " + constants.INVALID_OTP,
		})
		log.Println(constants.USER_NOT_FOUND+" or "+constants.INVALID_OTP, err.Error())
		return
	}

	if time.Now().After(otpExpiry) {
		config.WriteResponse(w, http.StatusUnauthorized, models.Response{
			Message: constants.OTP_EXPIRED,
		})
		log.Println(constants.OTP_EXPIRED)
		return
	}

	if strings.TrimSpace(req.OTP) != storedOTP {
		config.WriteResponse(w, http.StatusUnauthorized, models.Response{
			Message: constants.INVALID_OTP,
		})
		log.Println(constants.INVALID_OTP)
		return
	}

	insertUser := `INSERT INTO users (user_name, email, password, role, verified_at) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (email) DO UPDATE SET verified_at = $5`
	_, err = config.DB.Exec(insertUser, userName, req.Email, password, role, time.Now())

	if err != nil {
		config.WriteResponse(w, http.StatusInternalServerError, models.Response{
			Message: "Failed to finalize registration",
		})
		log.Printf("Failed to finalize registration. Error: %s\n", err)
		return
	}

	// Delete user from temp_users
	deleteUser := `DELETE FROM temp_users WHERE email = $1`
	_, err = config.DB.Exec(deleteUser, req.Email)
	if err != nil {
		log.Printf("Failed to delete temp user: %v\n", err)
	}

	config.WriteResponse(w, http.StatusCreated, models.Response{
		Message: "Email verified successfully. User Registered.",
	})
}

// VerifyOTPForgotPasswordHandler godoc
// @Summary Verify OTP for forgot password
// @Description Verify OTP and allow user to reset password
// @Tags Users
// @Accept  json
// @Produce  json
// @Param verifyRequest body models.VerifyRequest true "Verify OTP request"
// @Success 200 {object} models.ResponseWithEmail
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/forgot-password/verify-otp [post]
func VerifyOTPForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		config.WriteResponse(w, http.StatusBadRequest, models.Response{
			Message: constants.INVALID_REQUEST,
		})
		log.Printf(constants.INVALID_REQUEST+" Error: %s\n", err)
		return
	}

	ok := ValidateEmail(w, req.Email)
	if !ok {
		return
	}

	if len(req.OTP) != 6 {
		config.WriteResponse(w, http.StatusBadRequest, models.Response{
			Message: constants.INVALID_OTP,
		})
		log.Println(constants.INVALID_OTP)
		return
	}

	// store user details in variables
	var storedOTP string
	var otpExpiry time.Time

	userData := `SELECT otp, otp_expires FROM forgot_password WHERE email = $1`
	err := config.DB.QueryRow(userData, req.Email).Scan(&storedOTP, &otpExpiry)

	if err != nil {
		config.WriteResponse(w, http.StatusUnauthorized, models.Response{
			Message: constants.USER_NOT_FOUND + " or " + constants.INVALID_OTP,
		})
		log.Println(constants.USER_NOT_FOUND + " or " + constants.INVALID_OTP + err.Error())
		return
	}

	if time.Now().After(otpExpiry) {
		config.WriteResponse(w, http.StatusUnauthorized, models.Response{
			Message: constants.OTP_EXPIRED,
		})
		log.Println(constants.OTP_EXPIRED)
		return
	}

	if strings.TrimSpace(req.OTP) != storedOTP {
		config.WriteResponse(w, http.StatusUnauthorized, models.Response{
			Message: constants.INVALID_OTP,
		})
		log.Println(constants.INVALID_OTP)
		return
	}

	updateUser := `UPDATE forgot_password SET verified=$1 WHERE email=$2`
	_ = config.DB.QueryRow(updateUser, true, req.Email)

	config.WriteResponse(w, http.StatusOK, models.ResponseWithEmail{
		Message: "Email verified successfully",
		Email:   req.Email,
	})
}
