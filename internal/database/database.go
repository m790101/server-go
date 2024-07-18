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

// NewDB creates a new database connection
// and creates the database file if it doesn't exist
func NewDB(path string) (*DB, error) {

	mux := &sync.RWMutex{}
	db := &DB{
		path: path,
		mux:  mux,
	}
	return db, nil

}

// CreateChirp creates a new chirp and saves it to disk
func (db *DB) CreateChirp(body string) (Chirp, error) {
	data := Chirp{}
	return data, nil
}

// GetChirps returns all chirps in the database
func (db *DB) GetChirps() ([]Chirp, error) {
	data := []Chirp{}
	return data, nil
}

// ensureDB creates a new database file if it doesn't exist
func (db *DB) ensureDB() error {
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
	err := os.WriteFile("../../database.json", []byte("Hello, Gophers!"), 0666)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
