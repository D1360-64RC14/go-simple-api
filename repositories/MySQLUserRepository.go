package repositories

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/d1360-64rc14/simple-api/dtos"
	"github.com/d1360-64rc14/simple-api/interfaces"
	"github.com/d1360-64rc14/simple-api/utils"
)

// DefaultUserRepository implements UserRepository
var _ interfaces.UserRepository = (*MySQLUserRepository)(nil)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(database interfaces.Database) (interfaces.UserRepository, error) {
	repo := &MySQLUserRepository{
		db: database.DB(),
	}

	err := repo.createUserTableIfNotExist()
	if err != nil {
		return nil, err
	}

	return repo, nil
}

func (r MySQLUserRepository) Close() error {
	return r.db.Close()
}

func (r MySQLUserRepository) createUserTableIfNotExist() error {
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
// transaction not being started;
// transaction not being commited;
// query not being sucessfully executed;
// user just created not being found.
func (r MySQLUserRepository) CreateUser(user *dtos.UserWithHash) (*dtos.IdentifiedUser, *utils.ErrorCode) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, utils.NewErrorCode(http.StatusInternalServerError, err)
	}
	defer tx.Commit()

	_, err = tx.Exec(`
		INSERT INTO users(username, email, hash)
		VALUES (?, ?, ?);
	`, user.UserName, user.Email, user.Hash)
	if err != nil {
		return nil, utils.NewErrorCodeString(http.StatusConflict, "Email address already exist")
	}

	row := tx.QueryRow(`
		SELECT
			id
		FROM
			users
		WHERE
			username = ? AND
			email = ? AND
			hash = ?;
	`, user.UserName, user.Email, user.Hash)

	var id int
	err = row.Scan(&id)
	if err != nil {
		return nil, utils.NewErrorCode(http.StatusInternalServerError, err)
	}

	return &dtos.IdentifiedUser{
		ID:        id,
		UserModel: user.UserModel,
	}, nil
}

// SelectUserFromId returns the user with their id.
//
// Errors can be caused by:
// id not being found.
func (r MySQLUserRepository) SelectUserFromId(id int) (*dtos.IdentifiedUser, *utils.ErrorCode) {
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

	if row.Err() != nil {
		return nil, utils.NewErrorCode(http.StatusInternalServerError, row.Err())
	}

	user := new(dtos.IdentifiedUser)

	err := row.Scan(&user.ID, &user.UserName, &user.Email)
	if err != nil {
		return nil, utils.NewErrorCode(http.StatusNotFound, err)
	}

	return user, nil
}

// SelectUserFromEmail returns the identified user from their email.
//
// Errors can be caused by:
// query not being sucessfully executed;
// no rows being found.
func (r MySQLUserRepository) SelectUserFromEmail(email string) (*dtos.IdentifiedUser, *utils.ErrorCode) {
	row := r.db.QueryRow(`
		SELECT
			id,
			email,
			username
		FROM
			users
		WHERE
			email = ?
	`, email)

	if row.Err() != nil {
		return nil, utils.NewErrorCode(http.StatusInternalServerError, row.Err())
	}

	user := new(dtos.IdentifiedUser)

	err := row.Scan(&user.ID, &user.Email, &user.UserName)
	if err != nil {
		return nil, utils.NewErrorCode(http.StatusNotFound, err)
	}

	return user, nil
}

// SelectUserHashFromId returns the user password hash from database.
//
// Errors can be caused by:
// query not being sucessfully executed;
// id not being found.
func (r MySQLUserRepository) SelectUserHashFromId(id int) (string, *utils.ErrorCode) {
	row := r.db.QueryRow(`
		SELECT
			hash,
		FROM
			users
		WHERE
			id = ?;
	`, id)

	if row.Err() != nil {
		return "", utils.NewErrorCode(http.StatusInternalServerError, row.Err())
	}

	var hash string

	err := row.Scan(&hash)
	if err != nil {
		return "", utils.NewErrorCode(http.StatusNotFound, err)
	}

	return hash, nil
}

// SelectCompleteUserFromId reutrns all user info from database.
//
// Errors can be caused by:
// query not being sucessfully executed;
// id not being found.
func (r MySQLUserRepository) SelectCompleteUserFromId(id int) (*dtos.IdentifiedUserWithHash, *utils.ErrorCode) {
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

	if row.Err() != nil {
		return nil, utils.NewErrorCode(http.StatusInternalServerError, row.Err())
	}

	user := new(dtos.IdentifiedUserWithHash)

	err := row.Scan(&user.ID, &user.UserName, &user.Email, &user.Hash)
	if err != nil {
		return nil, utils.NewErrorCode(http.StatusNotFound, err)
	}

	return user, nil
}

// SelectAllUsers returns a list of all users from database.
//
// Errors can be caused by:
// user count query not being successfully executed;
// no rows being found at user count;
// row being read wrongly;
// no content being found at some row.
func (r MySQLUserRepository) SelectAllUsers() ([]*dtos.IdentifiedUser, *utils.ErrorCode) {
	row := r.db.QueryRow(`
		SELECT
			count(id)
		FROM users;
	`)

	if row.Err() != nil {
		return nil, utils.NewErrorCode(http.StatusInternalServerError, row.Err())
	}

	var userCount int

	err := row.Scan(&userCount)
	if err != nil {
		return nil, utils.NewErrorCode(http.StatusInternalServerError, err)
	}

	rows, err := r.db.Query(`
		SELECT
			id,
			username,
			email
		FROM
			users;
	`) // TODO: Add pagination
	if err != nil {
		return nil, utils.NewErrorCode(http.StatusInternalServerError, err)
	}
	defer rows.Close()

	users := make([]*dtos.IdentifiedUser, 0, userCount)

	for rows.Next() {
		user := new(dtos.IdentifiedUser)

		if rows.Err() != nil {
			return nil, utils.NewErrorCode(http.StatusInternalServerError, rows.Err())
		}

		err := rows.Scan(&user.ID, &user.UserName, &user.Email)
		if err != nil {
			return nil, utils.NewErrorCode(http.StatusNotFound, err)
		}

		users = append(users, user)
	}

	return users, nil
}

// RemoveUser removes an user from the database.
//
// Errors can be caused by:
// transaction not being started;
// transaction not being commited;
// more than 1 user being found;
// fail to get number of affected rows.
func (r MySQLUserRepository) RemoveUser(id int) *utils.ErrorCode {
	transaction, err := r.db.Begin()
	if err != nil {
		return utils.NewErrorCode(http.StatusInternalServerError, err)
	}

	result, err := transaction.Exec(`
		DELETE FROM
			users
		WHERE
			id = ?;
	`, id)
	if err != nil {
		return utils.NewErrorCode(http.StatusBadRequest, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.NewErrorCode(http.StatusInternalServerError, err)
	}

	if rowsAffected > 1 {
		transaction.Rollback()
		return utils.NewErrorCodeString(
			http.StatusConflict,
			fmt.Sprintf("There was %d users with id %d. User not removed.", rowsAffected, id),
		)
	}

	err = transaction.Commit()
	if err != nil {
		return utils.NewErrorCode(http.StatusInternalServerError, err)
	}

	return nil
}

// UserExist checks if an user with given id is present in the database.
//
// Errors can be caused by:
// query not being sucessfully executed;
// no rows being found.
func (r MySQLUserRepository) UserExist(id int) (bool, *utils.ErrorCode) {
	row := r.db.QueryRow(`
		SELECT EXISTS (
			SELECT
				1
			FROM
				users
			WHERE
				id = ?
		);
	`, id)
	if row.Err() != nil {
		return false, utils.NewErrorCode(http.StatusInternalServerError, row.Err())
	}

	var userExist bool

	err := row.Scan(&userExist)

	if err != nil {
		return false, utils.NewErrorCode(http.StatusInternalServerError, err)
	}

	return userExist, nil
}

func (r MySQLUserRepository) UpdateUsername(id int, newUsername string) *utils.ErrorCode {
	transaction, err := r.db.Begin()
	if err != nil {
		return utils.NewErrorCode(http.StatusInternalServerError, err)
	}

	result, err := transaction.Exec(`
		UPDATE
			users
		SET
			username = ?
		WHERE
			id = ?;
	`, id, newUsername)
	if err != nil {
		return utils.NewErrorCode(http.StatusBadRequest, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.NewErrorCode(http.StatusInternalServerError, err)
	}

	if rowsAffected > 1 {
		transaction.Rollback()
		return utils.NewErrorCodeString(
			http.StatusConflict,
			fmt.Sprintf("There was %d users with id %d. User not removed.", rowsAffected, id),
		)
	}

	err = transaction.Commit()
	if err != nil {
		return utils.NewErrorCode(http.StatusInternalServerError, err)
	}

	return nil
}
