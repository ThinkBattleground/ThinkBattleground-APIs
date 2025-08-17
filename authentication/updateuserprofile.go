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
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body models.Users true "User object"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /user/update-profile [put]
func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var user models.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		config.WriteResponse(w, http.StatusBadRequest, models.Response{
			Message: constants.INVALID_REQUEST,
		})
		log.Println(constants.INVALID_REQUEST + err.Error())
		return
	}

	updateUser := `UPDATE users SET user_name=$1 WHERE email=$2`
	_ = config.DB.QueryRow(updateUser, user.UserName, user.Email)

	config.WriteResponse(w, http.StatusOK, models.Response{
		Message: constants.UPDATE_PROFILE,
	})
	log.Printf("%s", constants.UPDATE_PROFILE)
}
