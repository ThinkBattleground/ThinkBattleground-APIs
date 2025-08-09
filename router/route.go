package router

import (
	"net/http"
	"thinkbattleground-apis/authentication"

	"thinkbattleground-apis/middleware"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "thinkbattleground-apis/docs"
)

func HandleRoute() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/user/register", authentication.RegisterUser).Methods("POST")
	r.HandleFunc("/user/verify-otp", authentication.VerifyOTPHandler).Methods("POST")
	r.HandleFunc("/user/login", authentication.LoginUser).Methods("POST")
	r.Handle("/user/logout", middleware.Auth(http.HandlerFunc(authentication.LogoutUser))).Methods("GET")
	r.HandleFunc("/user/forgot-password", authentication.ForgotPassword).Methods("POST")
	r.HandleFunc("/user/forgot-password/verify-otp", authentication.VerifyOTPForgotPasswordHandler).Methods("POST")
	r.HandleFunc("/user/forgot-password/reset-password", authentication.ResetPasswordAfterForgotPassword).Methods("PUT")
	r.Handle("/user/reset-password", middleware.Auth(http.HandlerFunc(authentication.ResetPassword))).Methods("PUT")
	r.HandleFunc("/user/update-profile", authentication.UpdateUserProfile).Methods("PUT")
	r.Handle("/user/role/{id}", middleware.Auth(http.HandlerFunc(authentication.ChangeUserRole))).Methods("PUT")
	r.Handle("/user/profile", middleware.Auth(http.HandlerFunc(authentication.GetUserProfileByCookie))).Methods("GET")
	r.Handle("/users", middleware.Auth(http.HandlerFunc(authentication.ListUsers))).Methods("GET")
	r.Handle("/user/{id}", middleware.Auth(http.HandlerFunc(authentication.GetUserById))).Methods("GET")
	r.Handle("/users?email={email}", middleware.Auth(http.HandlerFunc(authentication.FilterUserByEmail))).Methods("GET")
	r.Handle("/user/{id}", middleware.Auth(http.HandlerFunc(authentication.DeleteUsers))).Methods("DELETE")

	// Swagger route
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	return r
}
