package database

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

var ErrAlreadyExists = errors.New("already exists")
var ErrNotExist = errors.New("not exists")

// GetUsers returns all chirps in the database
func (db *DB) GetUsers() ([]User, error) {
	db.mux.RLock()
	data, err := os.ReadFile(db.path)
	if err != nil {
		log.Fatal(err)
	}
	dataDb := DBStructure{}
	res := []User{}
	json.Unmarshal(data, &dataDb)
	for _, data := range dataDb.Users {
		res = append(res, data)
	}

	db.mux.RUnlock()
	return res, nil
}

// GetUser returns all chirps in the database
func (db *DB) GetUser(id int) (User, error) {
	db.mux.RLock()
	dBStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
	}
	res, ok := dBStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}
	db.mux.RUnlock()
	return res, nil
}

func (db *DB) CreateUser(email string, hashPassword string, RefreshToken string) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
	}
	db.mux.Lock()

	newId := len(dbData.Users) + 1
	newUser := User{
		Id:          newId,
		Email:       email,
		Password:    hashPassword,
		IsChirpyRed: false,
	}

	dbData.Users[newId] = newUser
	db.mux.Unlock()
	errWrite := db.writeDB(dbData)
	if errWrite != nil {
		log.Fatal(err)
	}

	return newUser, nil
}

func (db *DB) UpdateUser(id int, email string, hashPassword string, isChirp bool) (User, error) {
	dBStructure, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
	}
	userModified, ok := dBStructure.Users[id]
	if !ok {
		return User{}, ErrNotExist
	}

	userModified.Email = email

	userModified.Password = hashPassword

	userModified.IsChirpyRed = isChirp

	dBStructure.Users[id] = userModified

	errWrite := db.writeDB(dBStructure)
	if errWrite != nil {
		log.Fatal(err)
	}

	return userModified, nil
}
