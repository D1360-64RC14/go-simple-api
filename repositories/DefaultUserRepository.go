package repositories

import (
	"database/sql"

	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/models"
	"golang.org/x/crypto/bcrypt"
)

type DefaultUserRepository struct {
	db *sql.DB
}

func NewDefaultUserRepository(database *sql.DB) (UserRepository, error) {
	repo := &DefaultUserRepository{
		db: database,
	}

	err := repo.createUserTableIfNotExist()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (r DefaultUserRepository) Close() error {
	return r.db.Close()
}

func (r DefaultUserRepository) createUserTableIfNotExist() error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()

	tx.Exec(`
		CREATE TABLE IF NOT EXISTS users(
			id       INTEGER      NOT NULL PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(50)  NOT NULL,
			email    VARCHAR(100) NOT NULL UNIQUE,
			hash     CHAR(72)     NOT NULL
		);
	`)

	return nil
}

func (r DefaultUserRepository) CreateUser(username, email, password string) (*dtos.IdentifiedUser, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO users(username, email, hash)
		VALUES (?, ?, ?);
	`, username, email, string(hash))
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow(`
		SELECT
			id
		FROM
			users
		WHERE
			email = ? AND
			username = ? AND
			hash = ?;
	`, email, username, string(hash))

	var id int
	row.Scan(&id)

	return &dtos.IdentifiedUser{
		ID: id,
		UserModel: models.UserModel{
			UserName: username,
			Email:    email,
		},
	}, nil
}

func (r DefaultUserRepository) SelectUserFromId(id int) (*dtos.IdentifiedUser, error) {
	row := r.db.QueryRow(`
		SELECT
			id,
			username,
			email
		FROM
			users
		WHERE
			id = ?;
	`, id)

	user := new(dtos.IdentifiedUser)

	err := row.Scan(&user.ID, &user.UserName, &user.Email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r DefaultUserRepository) SelectUserHashFromId(id int) (string, error) {
	row := r.db.QueryRow(`
		SELECT
			hash,
		FROM
			users
		WHERE
			id = ?;
	`, id)

	var hash string

	err := row.Scan(&hash)
	if err != nil {
		return "", err
	}

	return hash, nil
}

func (r DefaultUserRepository) SelectCompleteUserFromId(id int) (*dtos.IdentifiedUserWithHash, error) {
	row := r.db.QueryRow(`
		SELECT
			id,
			username,
			email,
			hash
		FROM
			users
		WHERE
			id = ?;
	`, id)

	user := new(dtos.IdentifiedUserWithHash)

	err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.Hash)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r DefaultUserRepository) SelectAllUsers() ([]*dtos.IdentifiedUser, error) {
	rows, err := r.db.Query(`
		SELECT
			id,
			username,
			email
		FROM
			users
	`)
	if err != nil {
		return nil, err
	}

	users := make([]*dtos.IdentifiedUser, 0)

	for rows.Next() {
		user := new(dtos.IdentifiedUser)

		err := rows.Scan(&user.ID, &user.UserName, &user.Email)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}
