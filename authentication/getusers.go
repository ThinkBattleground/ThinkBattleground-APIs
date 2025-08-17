package authentication

import (
	"log"
	"net/http"
	"thinkbattleground-apis/config"
	"thinkbattleground-apis/constants"
	"thinkbattleground-apis/models"

	"github.com/gorilla/mux"
)

// ListUsers godoc
// @Summary List all users
// @Description Get all users from the database
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {array} []models.UserGetResponse
// @Failure 400 {object} models.Response
// @Router /users [get]
func ListUsers(w http.ResponseWriter, r *http.Request) {
	ok := config.CheckAdmin(w, r)
	if !ok {
		return
	}
	userData := `SELECT id, user_name, email, role from users`
	rows, err := config.DB.Query(userData)
	if err != nil {
		config.WriteResponse(w, http.StatusBadRequest, models.Response{
			Message: constants.USER_NOT_FOUND,
		})
		log.Println(constants.USER_NOT_FOUND + err.Error())
		return
	}

	defer rows.Close()

	var usersList []models.UserGetResponse
	for rows.Next() {
		var user models.UserGetResponse
		if err := rows.Scan(&user.Id, &user.UserName, &user.Email, &user.Role); err != nil {
			log.Fatal(err)
		}
		usersList = append(usersList, user)
	}

	if len(usersList) == 0 {
		config.WriteResponse(w, http.StatusOK, models.Response{
			Message: constants.USER_NOT_FOUND,
		})
		return
	}
	config.WriteResponse(w, http.StatusOK, usersList)
}

// GetUserById godoc
// @Summary Get a user by ID
// @Description Get user by ID from the database
// @Tags Users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.UserGetResponse
// @Failure 400 {object} models.Response
// @Router /user/{id} [get]
func GetUserById(w http.ResponseWriter, r *http.Request) {
	ok := config.CheckAdmin(w, r)
	if !ok {
		return
	}
	var user models.UserGetResponse

	params := mux.Vars(r)
	id := params["id"]

	userData := `SELECT id, user_name, email, role from users WHERE id=$1`
	err := config.DB.QueryRow(userData, id).Scan(&user.Id, &user.UserName, &user.Email, &user.Role)
	if err != nil {
		config.WriteResponse(w, http.StatusBadRequest, models.Response{
			Message: constants.USER_NOT_FOUND,
		})
		log.Println(constants.USER_NOT_FOUND + err.Error())
		return
	}
	config.WriteResponse(w, http.StatusOK, user)
}

// FilterUserByEmail godoc
// @Summary Get a user by email
// @Description Get user by email from the database
// @Tags Users
// @Accept  json
// @Produce  json
// @Param email query string false "User Email (optional)"
// @Success 200 {object} []models.UserGetResponse
// @Failure 400 {object} models.Response
// @Router /users [get]
func FilterUserByEmail(w http.ResponseWriter, r *http.Request) {
	ok := config.CheckAdmin(w, r)
	if !ok {
		return
	}
	var user models.UserGetResponse

	params := mux.Vars(r)
	email := params["email"]

	userData := `SELECT id, user_name, email, role from users WHERE email=$1`
	err := config.DB.QueryRow(userData, email).Scan(&user.Id, &user.UserName, &user.Email, &user.Role)
	if err != nil {
		config.WriteResponse(w, http.StatusBadRequest, models.Response{
			Message: constants.USER_NOT_FOUND,
		})
		log.Println(constants.USER_NOT_FOUND + err.Error())
		return
	}
	config.WriteResponse(w, http.StatusOK, user)
}
