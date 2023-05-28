package repositories

import (
	"database/sql"

	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/d1360-64rc14/simple-api/models"
	"golang.org/x/crypto/bcrypt"
)

// DefaultUserRepository implements UserRepository
var _ interfaces.UserRepository = (*DefaultUserRepository)(nil)

type DefaultUserRepository struct {
	db *sql.DB
}

func NewDefaultUserRepository(database *sql.DB) (interfaces.UserRepository, error) {
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

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS users(
			id       INTEGER      NOT NULL PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(50)  NOT NULL,
			email    VARCHAR(100) NOT NULL UNIQUE,
			hash     CHAR(72)     NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	return nil
}

// CreateUser adds a new user to the database, returning an identified user.
//
// Errors can be caused by:
// transaction not beeing started;
// transaction not beeing commited;
// password having more than 72 characters;
// query not beeing sucessfully executed;
// user just created not beeing found.
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
	err = row.Scan(&id)
	if err != nil {
		return nil, err
	}

	return &dtos.IdentifiedUser{
		ID: id,
		UserModel: models.UserModel{
			UserName: username,
			Email:    email,
		},
	}, nil
}

// SelectUserFromId returns the user with their id.
//
// Errors can be caused by:
// id not beeing found.
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

// SelectUserHashFromId returns the user password hash from database.
//
// Errors can be caused by:
// id not beeing found.
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

// SelectCompleteUserFromId reutrns all user info from database.
//
// Errors can be caused by:
// id not beeing found.
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

// SelectAllUsers returns a list of all users from database.
//
// Errors can be caused by:
// query not beeing successfully executed;
// row beeing read wrongly.
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
	defer rows.Close()

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