package main

import (
	"fmt"
	"log"
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
		log.Fatalf("Error loading .env file")
	}

	err = db.Migrate()
	if err == nil {
		err := db.SeedUsers()
		if err != nil {
			fmt.Println("Could not seed users")
		}
	}
}

func main() {
	r := mux.NewRouter()
	finalhandler := http.HandlerFunc(user.GetProfile)

	r.Handle("/getProfile", middleware.Middleware((finalhandler))).Methods("GET")
	r.HandleFunc("/login", user.Login).Methods("PUT")

	fmt.Printf("Starting server at port 8000\n")
	http.ListenAndServe(":8000", r)

	removeExpiredKeys()
}

func removeExpiredKeys() {
	tickChan := time.NewTicker(tickerInterval).C

	for {
		<-tickChan
		fmt.Println("Now:", time.Now())
		db := db.DbConn()

		err := db.Where("expiration_time < ?", time.Now()).Delete(types.AuthKey{}).Error
		if err != nil {
			fmt.Println("Error deleting expired token: ", err)
		}
	}
}
