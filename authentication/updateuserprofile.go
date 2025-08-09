package authentication

import (
	"encoding/json"
	"log"
	"net/http"
	"thinkbattleground-apis/config"
	"thinkbattleground-apis/constants"
	"thinkbattleground-apis/models"
)

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Update user's profile information
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.Users true "User object"
// @Success 200 {object} models.Response
// @Failure 400 {object} map[string]string
// @Router /users/profile [put]
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var user models.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		config.WriteResponse(w, http.StatusBadRequest, constants.INVALID_REQUEST)
		log.Println(constants.INVALID_REQUEST + err.Error())
		return
	}

	updateUser := `UPDATE users SET first_name=$1, last_name=$2 WHERE email=$3`
	_ = config.DB.QueryRow(updateUser, user.FirstName, user.LastName, user.Email)

	resp := models.Response{
		Message: constants.UPDATE_PROFILE,
	}

	config.WriteResponse(w, http.StatusOK, resp)
	log.Printf(constants.PASSWORD_UPDATE)
}
