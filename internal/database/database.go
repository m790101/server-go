package database

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

type DB struct {
	path string
	mux  *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserRes struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type ParametersLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {

	mux := &sync.RWMutex{}
	db := &DB{
		path: path,
		mux:  mux,
	}
	db.ensureDB()
	return db, nil

}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mux.Lock()
	dbData, err := db.loadDB()
	if err != nil {
		log.Fatal(err)
	}

	newId := len(dbData.Chirps) + 1
	newChirp := Chirp{
		Id:   newId,
		Body: body,
	}
	dbData.Chirps[newId] = newChirp
	errWrite := db.writeDB(dbData)
	if errWrite != nil {
		log.Fatal(err)
	}

	db.mux.Unlock()
	return newChirp, nil
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

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	data, err := os.ReadFile(db.path)
	if err != nil {
		log.Fatal(err)
	}
	dataDb := DBStructure{}
	res := []Chirp{}
	json.Unmarshal(data, &dataDb)
	for _, data := range dataDb.Chirps {
		res = append(res, data)
	}

	db.mux.RUnlock()
	return res, nil
}

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

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.RLock()
	_, err := os.ReadFile(db.path)
	if err != nil {
		if err == os.ErrNotExist {
			_, errCreate := os.Create(db.path)
			initialData := DBStructure{}
			db.writeDB(initialData)
			if errCreate != nil {
				log.Fatal(err)
			}

		} else {
			log.Fatal(err)
		}
	}
	db.mux.RUnlock()
	return nil
}

// loadDB reads the database file into memory
func (db *DB) loadDB() (DBStructure, error) {
	db.mux.RLock()
	data, err := os.ReadFile(db.path)
	if err != nil {
		log.Fatal(err)
	}
	dbData := DBStructure{}
	json.Unmarshal(data, &dbData)

	db.mux.RUnlock()

	return dbData, nil

}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()

	dataJson, err := json.Marshal(dbStructure)

	errWrite := os.WriteFile(db.path, []byte(dataJson), 0644)
	if errWrite != nil {
		log.Fatal(err)
	}

	db.mux.Unlock()
	return nil
}
