package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"thinkbattleground-apis/config"
	"thinkbattleground-apis/constants"

	"github.com/golang-jwt/jwt"
)

var Claims jwt.MapClaims

func Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")

		if err != nil {
			if err != http.ErrNoCookie {
				config.WriteResponse(w, http.StatusUnauthorized, constants.USER_UNAUTHORIZED)
				log.Println(constants.USER_UNAUTHORIZED + err.Error())
				return
			}
			if cookie == nil {
				config.WriteResponse(w, http.StatusUnauthorized, "Error in Cookie: "+constants.USER_UNAUTHORIZED)
				log.Println(constants.USER_UNAUTHORIZED + err.Error())
				return
			}
			config.WriteResponse(w, http.StatusBadRequest, constants.BAD_REQUEST)
			log.Println(constants.BAD_REQUEST)
			return
		}

		tokenString := cookie.Value
		if err = config.LoadEnv(); err != nil {
			log.Println(constants.LOAD_ENV_ERROR)
			return
		}

		jwtKey := os.Getenv("JWTKEY")

		// parsing JWT string and storing the result in claims
		Claims = jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, Claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				// invalid signature
				config.WriteResponse(w, http.StatusUnauthorized, "Invalid Signature Error: "+constants.USER_UNAUTHORIZED)
				log.Println(constants.USER_UNAUTHORIZED + err.Error())
				return
			}
			config.WriteResponse(w, http.StatusBadRequest, "Error in parsing JWT: "+constants.BAD_REQUEST)
			log.Println("Error in parsing JWT: " + constants.BAD_REQUEST)
			return
		}

		if !token.Valid {
			config.WriteResponse(w, http.StatusUnauthorized, "Token Invalid: "+constants.USER_UNAUTHORIZED)
			log.Println("Token Invalid: " + constants.USER_UNAUTHORIZED)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			config.WriteResponse(w, http.StatusUnauthorized, constants.INVALID_TOKEN)
			log.Println(constants.INVALID_TOKEN)
			return
		}

		// Set user role in context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_data", claims)
		ctx = context.WithValue(ctx, "role", claims["role"])
		ctx = context.WithValue(ctx, "email", claims["email"])
		r = r.WithContext(ctx)

		// defining header's content-type
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
}
