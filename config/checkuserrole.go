package config

import (
	"log"
	"net/http"
	"strings"
)

func CheckAdmin(w http.ResponseWriter, r *http.Request) bool {
	resp := r.Context().Value("role").(string)
	if strings.Compare(resp, "admin") != 0 {
		WriteResponse(w, http.StatusBadRequest, "You do not have admin access")
		log.Println("You do not have admin access")
		return false
	}
	return true
}

func CheckFaculty(w http.ResponseWriter, r *http.Request) bool {
	resp := r.Context().Value("role").(string)
	if strings.Compare(resp, "faculty") != 0 {
		WriteResponse(w, http.StatusBadRequest, "You do not have faculty access")
		log.Println("You do not have faculty access")
		return false
	}
	return true
}
