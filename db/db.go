package db

import (
	"fmt"
	"log"
	"os"

	"github.com/shyams2012/Go_mini_project/seed"
	"github.com/shyams2012/Go_mini_project/types"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Mysql database connection
func DbConn() *gorm.DB {
	dbUserName := os.Getenv("MYSQL_USERNAME")
	dbPassword := os.Getenv("MYSQL_PASS")

	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/mini_project?charset=utf8mb4&parseTime=True&loc=Local", dbUserName, dbPassword)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	return db
}

// Migrate to DB
func Migrate() error {
	db := DbConn()
	return db.AutoMigrate(
		&types.User{},
		&types.AuthKey{},
	)
}

// Seeding of users
func SeedUsers() (err error) {
	db := DbConn()

	for _, seed := range seed.All() {
		if err = seed.Run(db); err != nil {
			log.Fatalf("Running seed '%s', failed with error: %s", seed.Name, err)
		}
	}
	return
}
