package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type user struct {
	gorm.Model

	Phone            string `gorm:"type:varchar(64);unique_index"`
	Email            string
	Name             string
	Surname          string
	Gender           string
	BirthDate        *time.Time
	RegistrationDate *time.Time

	Cars []car
}

func (user) TableName() string {
	return "users"
}

type carBrand struct {
	ID     uint `gorm:"primary_key"`
	NameEn string
}

func (carBrand) TableName() string {
	return "car_brands"
}

type carModel struct {
	ID      uint `gorm:"primary_key"`
	BrandID uint `gorm:"not null"`
	NameEn  string

	Brand carBrand `gorm:"ForeignKey:BrandID"`
}

func (carModel) TableName() string {
	return "car_models"
}

type car struct {
	ID      uint   `gorm:"primary_key"`
	UserID  uint   `gorm:"not null"`
	ModelID uint64 `gorm:"not null"`
	Color   string

	User  user     `gorm:"ForeignKey:UserID"`
	Model carModel `gorm:"ForeignKey:ModelID"`
}

func (car) TableName() string {
	return "cars"
}

func main() {
	db, err := gorm.Open("mysql", "golang_int_user:testpass@tcp(195.211.23.152:5435)/golang_intensive?parseTime=true")
	if err != nil {
		log.Fatalf("Can't connect to database: %s", err)
	}

	defer db.Close()

	db.LogMode(true)

	db.AutoMigrate(&user{})
	db.AutoMigrate(&car{})
	db.AutoMigrate(&carModel{})
	db.AutoMigrate(&carBrand{})

	db.Exec("DELETE FROM cars")
	db.Exec("DELETE FROM users")

	var users []uint
	for _, x := range [][]string{
		{"Vasya", "Pupkin", "79000000001"},
		{"Petya", "Ivanov", "79000000002"},
		{"Ivan", "Petrov", "79000000003"},
	} {
		u := user{
			Name:    x[0],
			Surname: x[1],
			Phone:   x[2],
		}

		if err := db.Create(&u).Error; err != nil {
			log.Fatalf("Can't store user: %s", err)
		}

		users = append(users, u.ID)
	}

	var nModels uint
	db.Model(&carModel{}).Count(&nModels)

	for i := 0; i < 10; i++ {
		u := rand.Int31n(int32(len(users)))
		model := rand.Int31n(int32(nModels))

		err := db.Create(&car{
			UserID:  users[u],
			ModelID: uint64(model),
			Color:   "#00ff00",
		}).Error

		if err != nil {
			log.Fatalf("Can't store car: %s", err)
		}
	}

	var selectedUser user
	selectedUser.ID = users[0]

	err = db.Model(&selectedUser).
		Preload("User").
		Preload("Model").
		Preload("Model.Brand").
		Related(&selectedUser.Cars).
		Error

	if err != nil {
		log.Fatalf("Can't get user cars: %s", err)
	}

	data, _ := json.MarshalIndent(selectedUser, "", "   ")
	log.Printf("%s", data)
}
