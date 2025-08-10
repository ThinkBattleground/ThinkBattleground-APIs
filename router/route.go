package router

import (
	"net/http"
	"thinkbattleground-apis/authentication"

	"thinkbattleground-apis/middleware"

	_ "thinkbattleground-apis/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func HandleRoute() *mux.Router {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/user/register", authentication.RegisterUser).Methods("POST", "OPTIONS")
	api.HandleFunc("/user/verify-otp", authentication.VerifyOTPHandler).Methods("POST")
	api.HandleFunc("/user/login", authentication.LoginUser).Methods("POST")
	api.Handle("/user/logout", middleware.Auth(http.HandlerFunc(authentication.LogoutUser))).Methods("GET")
	api.HandleFunc("/user/forgot-password", authentication.ForgotPassword).Methods("POST")
	api.HandleFunc("/user/forgot-password/verify-otp", authentication.VerifyOTPForgotPasswordHandler).Methods("POST")
	api.HandleFunc("/user/forgot-password/reset-password", authentication.ResetPasswordAfterForgotPassword).Methods("PUT")
	api.Handle("/user/reset-password", middleware.Auth(http.HandlerFunc(authentication.ResetPassword))).Methods("PUT")
	api.HandleFunc("/user/update-profile", authentication.UpdateUserProfile).Methods("PUT")
	api.Handle("/user/role/{id}", middleware.Auth(http.HandlerFunc(authentication.ChangeUserRole))).Methods("PUT")
	api.Handle("/user/profile", middleware.Auth(http.HandlerFunc(authentication.GetUserProfileByCookie))).Methods("GET")
	api.Handle("/users", middleware.Auth(http.HandlerFunc(authentication.ListUsers))).Methods("GET")
	api.Handle("/user/{id}", middleware.Auth(http.HandlerFunc(authentication.GetUserById))).Methods("GET")
	api.Handle("/users?email={email}", middleware.Auth(http.HandlerFunc(authentication.FilterUserByEmail))).Methods("GET")
	api.Handle("/user/{id}", middleware.Auth(http.HandlerFunc(authentication.DeleteUsers))).Methods("DELETE")

	// Swagger route (no prefix)
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return r
}
