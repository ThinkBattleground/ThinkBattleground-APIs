package authentication

import (
	"fmt"
	"net/http"
	"thinkbattleground-apis/config"
	"thinkbattleground-apis/models"
	"time"
)

// LogoutUser godoc
// @Summary Logout user
// @Description Logout user and clear JWT token
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Response
// @Router /user/logout [post]
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	email := r.Context().Value("email").(string)
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Path:    "/",
		Expires: time.Now(),
	})
	config.WriteResponse(w, http.StatusOK, models.Response{
		Message: fmt.Sprintf("User %s logged out successfully!", email),
	})
}
