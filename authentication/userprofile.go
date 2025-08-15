package authentication

import (
	"net/http"
	"thinkbattleground-apis/config"
)

// GetUserProfileByCookie godoc
// @Summary Get user profile by cookie
// @Description Get user profile using cookie
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]interface{}
// @Router /user/profile [get]
func GetUserProfileByCookie(w http.ResponseWriter, r *http.Request) {
	resp := r.Context().Value("user_data")
	config.WriteResponse(w, http.StatusOK, resp)
}
