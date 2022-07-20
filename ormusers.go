package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

type User struct {
	gorm.Model
	Name  string
	Email string
}

func connectDB() (*gorm.DB, error) {
	dsn := "host=0.0.0.0 user=postgres password=password dbname=ormdb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return db, err
}

func InitiateMigration() {
	db, err := connectDB()
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}
	defer db.Close()

	db.AutoMigrate(&User{})
}

func AllUsers(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}
	defer db.Close()

	var users []User
	db.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func NewUser(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}
	defer db.Close()

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err)
	}

	db.Create(&user)

	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	var user User
	db.Where("id = ?", id).Find(&user)

	var reqJSON map[string]string
	err = json.NewDecoder(r.Body).Decode(&reqJSON)
	if name, ok := reqJSON["name"]; ok {
		user.Name = name
	}

	if email, ok := reqJSON["email"]; ok {
		user.Email = email
	}

	db.Save(&user)

	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to DB")
	}
	defer db.Close()

	vars := mux.Vars(r)
	id := vars["id"]

	var user User
	db.Where("id = ?", id).Find(&user)

	db.Delete(&user)

	fmt.Fprintf(w, "User %s is deleted.", user.Name)
}
