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
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
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
		log.Printf("Failed to delete data from user: %v\n", err)
	}

	resp := models.Response{
		Message: "User deleted successfully",
	}

	config.WriteResponse(w, http.StatusOK, resp)
}
