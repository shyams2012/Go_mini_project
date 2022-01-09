package middleware

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/shyams2012/Go_mini_project/auth"
	"github.com/shyams2012/Go_mini_project/user"
)

//middleware for api authorization
// func Middleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		_, err := auth.ValidateToken(r.Header.Get("Authorization"))

// 		if err != nil {
// 			json.NewEncoder(w).Encode(fmt.Sprintf("Unauthorised. Error: %w", err))
// 		} else {
// 			next.ServeHTTP(w, r)
// 		}
// 	})
// }
// assign the secret key to key variable on program's first run
// func init() {
// 	key := make([]byte, 64)
// 	_, err := rand.Read(key)
// 	if err != nil {
// 		fmt.Println("err is", err)
// 	}

// 	fmt.Println(key)
// }

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeaderParts := strings.Split(r.Header.Get("Authorization"), "Bearer ")

		if len(authHeaderParts) != 2 {
			fmt.Println("errrr")

			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {

			jwtToken := authHeaderParts[1]

			// Get claims without parsing JWT
			unsafeClaims, err := auth.GetUnsafeClaims(jwtToken)
			if err != nil {
				log.Fatalf("failed to parse to jwt: %v", err)
			}

			//fmt.Println("Unsafe User's email:", unsafeClaims.Email)
			keys := user.GetPublicKeys(unsafeClaims.Email)

			//fmt.Println("keys", keys)
			//fmt.Println("keys", timedata)

			var claims *auth.UserClaims

			for _, key := range keys {

				var publicKey *rsa.PublicKey
				if err := json.Unmarshal([]byte(key.PublicKey), &publicKey); err != nil {
					fmt.Println(err)
				}

				//Validate JWT with JWK
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
