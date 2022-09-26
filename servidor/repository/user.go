package repository

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

/*
#     #  #####  ####### ######
#     # #     # #       #     #
#     # #       #       #     #
#     #  #####  #####   ######
#     #       # #       #   #
#     # #     # #       #    #
 #####   #####  ####### #     #
*/

func (r *Repository) CreateUser(user User) error {

	_, erri := r.db.Exec("INSERT INTO Users SET Admin = ?, Banned = ?, Avatar = ?, Verified = ?, StatusID = ?, Username = ?, Pass = ?, Email = ?, EmailCode = ?", user.Admin, user.Banned, user.Avatar, user.Verified, user.StatusID, user.Username, user.Pass, user.Email, user.EmailCode)
	if erri != nil {
		return errors.New("Error adding user")
	}

	return nil
}

func (r *Repository) DeleteUser(ID int) (bool, error) {

	result, erri := r.db.Exec("DELETE FROM Users WHERE ID = ?", ID)

	if erri != nil {
		return false, errors.New("Error deleting user")
	}

	rows, errr := result.RowsAffected()

	if errr != nil {
		return false, errors.New("Internal server error")
	}

	if rows == 0 {
		return false, errors.New("User does not exist")
	}

	return true, nil
}

func (r *Repository) GetUser(ID int) (User, error) {

	var user = User{}

	errs := r.db.QueryRow(`SELECT * FROM Users WHERE ID = ?`, ID).Scan(&user.ID, &user.Username, &user.Pass, &user.Email, &user.EmailCode, &user.Avatar, &user.Admin, &user.Banned, &user.Verified, &user.StatusID)
	if errs == sql.ErrNoRows {
		return User{}, errors.New("User does not exist")
	} else if errs != nil {
		return User{}, errors.New("Error querying user")
	}

	return user, nil
}

func (r *Repository) GetUserByCode(code string) (User, error) {

	var user = User{}

	errs := r.db.QueryRow(`SELECT * FROM Users WHERE EmailCode = ?`, code).Scan(&user.ID, &user.Username, &user.Pass, &user.Email, &user.EmailCode, &user.Avatar, &user.Admin, &user.Banned, &user.Verified, &user.StatusID)

	if errs == sql.ErrNoRows {
		return User{}, errors.New("User does not exist")
	} else if errs != nil {
		return User{}, errors.New("Error querying user")
	}

	return user, nil
}

func (r *Repository) GetUserByEmail(email string) (User, error) {

	var user = User{}

	errs := r.db.QueryRow(`SELECT * FROM Users WHERE Email = ?`, email).Scan(&user.ID, &user.Username, &user.Pass, &user.Email, &user.EmailCode, &user.Avatar, &user.Admin, &user.Banned, &user.Verified, &user.StatusID)
	if errs == sql.ErrNoRows {
		return User{}, errors.New("User does not exist")
	} else if errs != nil {
		return User{}, errors.New("Error querying user")
	}

	return user, nil
}

func (r *Repository) UpdateUser(ID int, user User) error {

	result, erru := r.db.Exec(`UPDATE Users SET Username = ?,Pass = ?,Email = ?,EmailCode = ?,Avatar = ?,Admin = ?,Banned = ?,Verified = ?,StatusID = ? WHERE ID = ?`, user.Username, user.Pass, user.Email, user.EmailCode, user.Avatar, user.Admin, user.Banned, user.Verified, user.StatusID, ID)
	if erru != nil {
		return errors.New("Error updating user")
	}

	rows, errr := result.RowsAffected()

	if errr != nil {
		return errors.New("Internal server error")
	}

	if rows == 0 {
		return errors.New("User does not exist")
	}

	return nil
}

func (r *Repository) ExistsUserByEmail(email string) bool {
	var user = User{}

	err := r.db.QueryRow("SELECT ID FROM Users WHERE Email = ?", email).Scan(&user.ID)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	return true
}

func (r *Repository) ExistsUser(ID int) bool {
	var user = User{}

	err := r.db.QueryRow("SELECT ID FROM Users WHERE ID = ?", ID).Scan(&user.ID)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	return true
}
