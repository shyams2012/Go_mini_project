package main

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/shyams2012/Go_mini_project/db"
	"github.com/shyams2012/Go_mini_project/middleware"
	"github.com/shyams2012/Go_mini_project/types"
	"github.com/shyams2012/Go_mini_project/user"
)

const (
	tickerInterval = time.Second * 10
)

func init() {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error getting env file. Error :", err)
	}

	// Migrate to DB
	err = db.Migrate()
	if err == nil {
		// Seed users
		err := db.SeedUsers()
		if err != nil {
			fmt.Println("Could not seed users")
		}
	}
}

func main() {
	r := mux.NewRouter()

	// Handle get user profile
	finalhandler := http.HandlerFunc(user.GetProfile)
	r.Handle("/getProfile", middleware.Middleware((finalhandler))).Methods("GET")
	// Handle get user login
	r.HandleFunc("/login", user.Login).Methods("PUT")

	// Remove expired keys in background
	go removeExpiredKeys()

	fmt.Printf("Starting server at port 8000\n")

	// Listen to server
	http.ListenAndServe(":8000", r)
}

//Removes expired keys from database
func removeExpiredKeys() {
	fmt.Println("Start removing expired keys in background")

	tickChan := time.NewTicker(tickerInterval).C

	//Infinite loop to remove expired keys using NewTicker function
	for {
		<-tickChan
		db := db.DbConn()
		err := db.Where("expiration_time < ?", time.Now()).Delete(types.AuthKey{}).Error
		if err != nil {
			fmt.Println("Error deleting expired token.Error:", err)
		}
	}
}
