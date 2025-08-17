package authentication

import (
	"encoding/json"
	"log"
	"net/http"
	"thinkbattleground-apis/config"
	"thinkbattleground-apis/constants"
	"thinkbattleground-apis/models"

	"github.com/gorilla/mux"
)

// ChangeUserRole godoc
// @Summary Change a user's role
// @Description Change the role of a user by admin
// @Tags Users
// @Accept  json
// @Produce  json
// @Param user body models.Users true "User object"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /user/role/{id} [put]
func ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	ok := config.CheckAdmin(w, r)
	if !ok {
		return
	}

	var user models.Users
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		config.WriteResponse(w, http.StatusBadRequest, models.Response{
			Message: constants.INVALID_REQUEST,
		})
		log.Printf(constants.INVALID_REQUEST+" Error: %s\n", err)
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	updateUser := `UPDATE users SET role=$1 WHERE id=$2`
	_ = config.DB.QueryRow(updateUser, user.Role, id)

	config.WriteResponse(w, http.StatusOK, models.Response{
		Message: "User updated successfully",
	})
}
