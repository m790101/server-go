package database

import (
	"encoding/json"
	"log"
	"os"
)

// GetChirps returns all chirps in the database
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

func (db *DB) CreateUser(body ParametersLogin) (User, error) {
	dbData, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
	}
	db.mux.Lock()

	newId := len(dbData.Users) + 1
	newUser := User{
		Id:       newId,
		Email:    body.Email,
		Password: body.Password,
	}

	dbData.Users[newId] = newUser
	db.mux.Unlock()
	errWrite := db.writeDB(dbData)
	if errWrite != nil {
		log.Fatal(err)
	}

	return newUser, nil
}
