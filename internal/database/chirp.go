package database

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

func (db *DB) CreateChirp(body string, userId int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	uuid := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:       uuid,
		Body:     body,
		AuthorId: userId,
	}
	dbStructure.Chirps[uuid] = chirp

	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	chirp, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, ErrNotExist
	}

	return chirp, nil
}

func (db *DB) DeleteChirp(id int) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}

	// chirp, ok := dbStructure.Chirps[id]

	delete(dbStructure.Chirps, id)

	return nil
}
