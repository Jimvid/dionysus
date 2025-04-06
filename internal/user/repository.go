package user

import (
	"database/sql"
	"fmt"

	"github.com/jimvid/dionysus/internal/model"
)

type UserRepository struct {
	db *sql.DB
}

const USER_TABLE_NAME = "userTable"

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (u UserRepository) DoesUserExist(username string) (bool, error) {
	const query = "SELECT 1 FROM users WHERE username = ? LIMIT 1"

	var exists int
	err := u.db.QueryRow(query, username).Scan(&exists)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (u UserRepository) InsertUser(user model.User) error {
	query := "INSERT INTO users (username, password, email) VALUES (?, ?, ?)"

	_, err := u.db.Exec(query, user.Username, user.PasswordHash)
	if err != nil {
		return err
	}

	return nil
}

func (u UserRepository) GetUser(username string) (model.User, error) {
	var user model.User

	query := "SELECT username, email FROM users WHERE username = ?"

	err := u.db.QueryRow(query, username).Scan(&user.Username, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("user not found")
		}
		return user, err
	}

	return user, nil
}
