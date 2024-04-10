package database

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	path string
	mu   *sync.RWMutex
}

type DBStructure struct {
    Users map[int]User `json:"users"`
	Chirps map[int]Chirp `json:"chirps"`
}

type User struct {
    ID   int    `json:"id"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
    ID   int    `json:"id"`
	Email string `json:"email"`
}
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func NewDB(path string) (*DB, error) {
	dbg := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()
	
	if *dbg {
        os.Remove(path)
    }
	
	db := &DB{
		path: path,
		mu:   &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) CreateUser(email string, password []byte) (UserResponse, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return UserResponse{}, err
	}

	for _, user := range dbStructure.Users {
        if user.Email == email {
			return UserResponse{}, errors.New("couldn't create user")
        }
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, err
	}

	id := len(dbStructure.Users) + 1
	user := User{
		ID:   id,
		Email: email,
		Password: string(hashedPassword),
	}
	dbStructure.Users[id] = user

	err = db.writeDB(dbStructure)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{
		ID: user.ID,
		Email: user.Email,
	}, nil
}

func (db *DB) LoginUser(email string, password []byte) (UserResponse, error) {
	dbStructure, err := db.loadDB()
    if err != nil {
        return UserResponse{}, err
    }

    for _, user := range dbStructure.Users {
        if user.Email != email {
			continue;
        }

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			return UserResponse{}, err
		} else {
			return UserResponse{
				ID: user.ID,
				Email: user.Email,
			}, nil
		}
    }

    return UserResponse{}, nil
}

func (db *DB) GetUsers() ([]UserResponse, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	users := make([]UserResponse, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, UserResponse{
			ID: user.ID,
			Email: user.Email,
		})
	}

	return users, nil
}

func (db *DB) GetUserByID(userID int) (UserResponse, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return UserResponse{}, err
    }

    for _, user := range dbStructure.Users {
        if user.ID == userID {
            return UserResponse{
				ID: user.ID,
				Email: user.Email,
			}, nil
        }
    }

    return UserResponse{}, nil
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(dbStructure.Chirps) + 1
	chirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[id] = chirp

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

func (db *DB) GetChirpByID(chirpID int) (Chirp, error) {
    dbStructure, err := db.loadDB()
    if err != nil {
        return Chirp{}, err
    }

    for _, chirp := range dbStructure.Chirps {
        if chirp.ID == chirpID {
            return chirp, nil
        }
    }

    return Chirp{}, nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
        Users: map[int]User{},
		Chirps: map[int]Chirp{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}

	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}