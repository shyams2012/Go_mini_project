package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/shyams2012/Go_mini_project/auth"
	"github.com/shyams2012/Go_mini_project/user"
)

//Middleware for authorization
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeaderParts := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeaderParts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			jwtToken := authHeaderParts[1]
			unsafeClaims, err := auth.GetUnsafeClaims(jwtToken)
			if err != nil {
				fmt.Println("Error getting unsafeClaims. Error :", err)
			}

			keys := user.GetPublicKeys(unsafeClaims.Email)

			var claims *auth.UserClaims
			for _, key := range keys {
				var publicKey *rsa.PublicKey
				if err := json.Unmarshal([]byte(key.PublicKey), &publicKey); err != nil {
					fmt.Println("Error getting publicKeys. Error :", err)
				}
				claims, err = auth.ParseJWT(jwtToken, publicKey)
				if err == nil {
					break
				}
			}

			if claims == nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			}

			userEmail := claims.Email
			ctx := context.WithValue(r.Context(), "email", userEmail)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
