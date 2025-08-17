package authentication

import (
	"log"
	"net/http"
	"thinkbattleground-apis/config"
	"thinkbattleground-apis/models"

	"github.com/gorilla/mux"
)

// DeleteUsers godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /user/{id} [delete]
func DeleteUsers(w http.ResponseWriter, r *http.Request) {
	ok := config.CheckAdmin(w, r)
	if !ok {
		return
	}

	params := mux.Vars(r)
	id := params["id"]

	deleteUser := `DELETE FROM users WHERE id = $1`
	_, err := config.DB.Exec(deleteUser, id)
	if err != nil {
		config.WriteResponse(w, http.StatusInternalServerError, models.Response{
			Message: "Failed to delete user",
		})
		log.Printf("Failed to delete data from user: %v\n", err)
		return
	}

	config.WriteResponse(w, http.StatusOK, models.Response{
		Message: "User deleted successfully",
	})
}
