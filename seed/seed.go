package seed

import (
	"fmt"
	"log"

	"github.com/shyams2012/Go_mini_project/types"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, email, password, name, location string) error {
	var users = []types.User{}

	if err := db.Where("email = ?", email).Find(&users).Error; err != nil {
		log.Fatalln(err)
		fmt.Println(err)
	}

	if len(users) > 0 {
		return nil
	}

	return db.Create(&types.User{
		Email:    email,
		Password: password,
		Name:     name,
		Location: location,
	}).Error
}

func All() []types.Seed {
	hashedPassword1, _ := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.DefaultCost)
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("password2"), bcrypt.DefaultCost)

	return []types.Seed{
		types.Seed{
			Name: "CreateShyam",
			Run: func(db *gorm.DB) error {
				return CreateUser(db, "shyams2012@gmail.com", string(hashedPassword1), "Shyam", "Ktm")
			},
		},
		types.Seed{
			Name: "CreateAjay",
			Run: func(db *gorm.DB) error {
				return CreateUser(db, "ajay@gmail.com", string(hashedPassword2), "Ajay", "Ktm")
			},
		},
	}
}
