package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func exec(q string, args ...interface{}) sql.Result {
	logArgs := append([]interface{}{"Trying to execute:", q}, args...)
	log.Println(logArgs...)

	result, err := db.Exec(q, args...)
	if err != nil {
		log.Fatalf("Can't execute query %q (%+v): %s", q, args, err)
	}

	return result
}

func query(q string, args ...interface{}) *sql.Rows {
	logArgs := append([]interface{}{"Trying to execute:", q}, args...)
	log.Println(logArgs...)

	result, err := db.Query(q, args...)
	if err != nil {
		log.Fatalf("Can't execute query %q (%+v): %s", q, args, err)
	}

	return result
}

func init() {
	database, err := sql.Open("mysql", "golang_int_user:testpass@tcp(195.211.23.152:5435)/golang_intensive")
	if err != nil {
		log.Fatalf("Can't connect to database: %s", err)
	}

	db = database

	exec("DELETE FROM cars")
	exec("DELETE FROM users")
}

type user struct {
	ID               uint64
	Phone            string
	Email            string
	Name             string
	Surname          string
	Gender           string
	BirthDate        time.Time
	RegistrationDate time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type carBrand struct {
	ID     uint64
	NameEn string
}

type carModel struct {
	ID      uint64
	BrandID uint64
	NameEn  string

	carBrand *carBrand
}

type car struct {
	ID      uint64
	UserID  uint64
	ModelID uint64
	Color   string

	carModel *carModel
	user     *user
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func addUser(u user) {
	exec(fmt.Sprintf(`INSERT INTO users (
		phone,
		email,
		name,
		surname,
		gender,
		birth_date,
		registration_date,
		created_at,
		updated_at
	) values (
		'%s',
		'%s',
		'%s',
		'%s',
		'%s',
		'%s',
		'%s',
		'%s',
		'%s'
	)`,
		u.Phone,
		u.Email,
		u.Name,
		u.Surname,
		u.Gender,
		formatTime(u.BirthDate),
		formatTime(u.RegistrationDate),
		formatTime(u.CreatedAt),
		formatTime(u.UpdatedAt),
	))
}

func addUser2(u *user) {
	ret := exec(`INSERT INTO users (
		phone,
		email,
		name,
		surname,
		gender,
		birth_date,
		registration_date,
		created_at,
		updated_at
	) values (
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	)`,
		u.Phone,
		u.Email,
		u.Name,
		u.Surname,
		u.Gender,
		u.BirthDate,
		formatTime(u.RegistrationDate),
		formatTime(u.CreatedAt),
		formatTime(u.UpdatedAt),
	)

	id, err := ret.LastInsertId()
	if err != nil {
		log.Fatalf("Can't get last inserted id: %s", err)
	}

	u.ID = uint64(id)
}

func addCar(c car) {
	exec(`INSERT INTO cars (
		user_id,
		color,
		model_id
	) values (
		?,
		?,
		?
	)`, c.UserID, c.Color, c.ModelID)
}

func getUserCars(userID uint64) []car {
	rows := query(`SELECT
		user.id,
		user.phone,
		user.email,
		user.name,
		user.surname,
		user.gender,
		user.birth_date,
		user.registration_date,
		user.created_at,
		user.updated_at,

		car_brand.id,
		car_brand.name_en,

		car_model.id,
		car_model.brand_id,
		car_model.name_en,

		car.id,
		car.user_id,
		car.model_id,
		car.color
	FROM cars car
	JOIN users user ON (user.id = car.user_id)
	JOIN car_models car_model ON (car.model_id = car_model.id)
	JOIN car_brands car_brand ON (car_model.brand_id = car_brand.id)
	WHERE user.id = ?`, userID)

	var cars []car

	defer rows.Close()
	for rows.Next() {
		var birthDate, registrationDate, createdAt, updatedAt string // TODO: convert to timestamps
		c := car{
			user: &user{},
			carModel: &carModel{
				carBrand: &carBrand{},
			},
		}

		err := rows.Scan(
			&c.user.ID,
			&c.user.Phone,
			&c.user.Email,
			&c.user.Name,
			&c.user.Surname,
			&c.user.Gender,
			&birthDate,
			&registrationDate,
			&createdAt,
			&updatedAt,

			&c.carModel.carBrand.ID,
			&c.carModel.carBrand.NameEn,

			&c.carModel.ID,
			&c.carModel.BrandID,
			&c.carModel.NameEn,

			&c.ID,
			&c.UserID,
			&c.ModelID,
			&c.Color,
		)

		if err != nil {
			log.Fatalf("Can't read car for user %d: %s", userID, err)
		}

		cars = append(cars, c)
	}

	return cars
}

func main() {
	addUser(user{
		Name:    "Vasya",
		Surname: "Pupkin",
	})

	newUser := user{
		Name:    "Ivan",
		Surname: "Popov",
		Phone:   "79990002233",
	}

	addUser2(&newUser)
	log.Printf("User with ID %d successfully inserted\n", newUser.ID)

	addCar(car{
		UserID:  newUser.ID,
		ModelID: 10,
		Color:   "#ff00ff",
	})

	cars := getUserCars(newUser.ID)
	log.Printf("%+v", cars)
}
