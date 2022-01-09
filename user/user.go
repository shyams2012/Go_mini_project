package user

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/shyams2012/Go_mini_project/auth"
	"github.com/shyams2012/Go_mini_project/db"
	"github.com/shyams2012/Go_mini_project/types"

	"golang.org/x/crypto/bcrypt"
)

// Token jwt Standard Claim Object
type Token struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var Expire_time = time.Now().Add(5 * time.Minute)

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		db := db.DbConn()
		var result types.User
		body, _ := ioutil.ReadAll(r.Body)
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println(err)
		}
		email := result.Email
		password := result.Password
		fmt.Println("result.Email", result.Email)

		user := types.User{}

		db.Where("Email = ?", email).Find(&user)

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return
		}

		// Create a token object
		var tokenObj = Token{
			Email: email,
			StandardClaims: jwt.StandardClaims{
				// Enter expiration in minutes
				ExpiresAt: Expire_time.Unix(),
			},
		}

		rsaKey, err := auth.NewRSAKey()

		if err != nil {
			log.Fatalf("failed to generate rsa key: %v", err)
		}

		// Generate JWT
		jwt, err := auth.NewJWT(tokenObj.Email, rsaKey, time.Now().Add(auth.JWTAddExpiry))
		if err != nil {
			log.Fatalf("failed to generate to jwt: %v", err)
		}

		// Add keys to DB
		privateKey, _ := json.Marshal(rsaKey)
		privateString := string(privateKey)
		publicKey, _ := json.Marshal(rsaKey.PublicKey)
		publicString := string(publicKey)

		SavePrivateKey_PublicKey(privateString, publicString, email)
		json.NewEncoder(w).Encode(jwt)
	}
}

func SavePrivateKey_PublicKey(privateString, publicString, email string) {
	db := db.DbConn()

	authKeys := &types.AuthKey{Email: email, PrivateKey: privateString, PublicKey: publicString, ExpirationTime: Expire_time}

	if err := db.Create(authKeys).Error; err != nil {
		fmt.Println("err: found", err)
	}
}

func GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		db := db.DbConn()
		user := types.User{}
		email := r.Context().Value("email")

		if err := db.Table("users").Select("email", "name", "location").Where("email = ?", email).Scan(&user).Error; err != nil {
			fmt.Print("Error getting user profile. Error :", err)
		}

		json.NewEncoder(w).Encode(user)
	}
}

func GetPublicKeys(email string) []types.AuthKey {
	publicKeys := make([]types.AuthKey, 0)
	publicKey := new(types.AuthKey)

	db := db.DbConn()
	rows, _ := db.Table("auth_Keys").Select("public_key").Where("email = ?", email).Rows()
	defer rows.Close()

	for rows.Next() {
		err := db.ScanRows(rows, &publicKey)
		if err != nil {
			fmt.Println("No public keys found")
		}
		publicKeys = append(publicKeys, *publicKey)
	}

	return publicKeys
}
