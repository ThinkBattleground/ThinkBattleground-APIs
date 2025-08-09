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

// ResetPassword godoc
// @Summary Reset password
// @Description Reset user password
// @Tags users
// @Accept  json
// @Produce  json
// @Param data body models.ResetPassword true "Reset password data"
// @Success 200 {object} models.ResponseWithEmail
// @Failure 400 {object} map[string]string
// @Router /users/reset-password [post]
func ResetPassword(w http.ResponseWriter, r *http.Request) {
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

	// encrypt password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)

	updateUser := `UPDATE users SET password=$1, verified_at=$2 WHERE email=$3`
	_ = config.DB.QueryRow(updateUser, hashedPassword, time.Now(), data.Email)

	resp := models.ResponseWithEmail{
		Message: constants.PASSWORD_UPDATE,
		Email:   data.Email,
	}

	config.WriteResponse(w, http.StatusOK, resp)
	log.Printf(constants.PASSWORD_UPDATE)
}
