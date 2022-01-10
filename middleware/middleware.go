package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/shyams2012/Go_mini_project/auth"
	"github.com/shyams2012/Go_mini_project/types"
	"github.com/shyams2012/Go_mini_project/user"
)

// Middleware for authorization
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeaderParts := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		// Check if token is not malformed
		if len(authHeaderParts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			// Take JWT token part
			jwtToken := authHeaderParts[1]
			unsafeClaims, err := auth.GetUnsafeClaims(jwtToken)
			if err != nil {
				fmt.Println("Error getting unsafeClaims. Error :", err)
			}

			// Get user's public key set
			keys := user.GetPublicKeys(unsafeClaims.Email)

			// Iterate over keys to find correct public key to parse JWT token
			var claims *types.UserClaims
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
			// Set user email in context
			ctx := context.WithValue(r.Context(), "email", claims.Email)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
