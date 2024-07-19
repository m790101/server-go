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

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

type ResponseType struct {
	Id   int    `json:"id"`
	Body string `json:"cleaned_body"`
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

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	db.mux.RLock()
	data, err := os.ReadFile("../../database.json")
	if err != nil {
		log.Fatal(err)
	}
	dataRes := []Chirp{}

	json.Unmarshal(data, &dataRes)

	db.mux.RUnlock()
	return dataRes, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
	db.mux.RLock()
	_, err := os.ReadFile("../../database.json")
	if err != nil {
		if err == os.ErrNotExist {
			_, errCreate := os.Create("../../database.json")
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
	data, err := os.ReadFile("../../database.json")
	if err != nil {
		log.Fatal(err)
	}
	dataJson := []Chirp{}

	json.Unmarshal(data, &dataJson)
	dbData := DBStructure{}
	for _, chirp := range dataJson {
		dbData.Chirps[chirp.Id] = chirp
	}

	db.mux.RUnlock()

	return dbData, nil

}

// writeDB writes the database file to disk
func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mux.Lock()
	data := []Chirp{}

	for _, chirp := range dbStructure.Chirps {
		data = append(data, chirp)
	}

	dataJson, err := json.Marshal(data)

	errWrite := os.WriteFile("../../database.json", []byte(dataJson), 0666)
	if errWrite != nil {
		log.Fatal(err)
	}

	db.mux.Unlock()
	return nil
}
