package repository

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

/*
######  #       #######  #####  #    #
#     # #       #     # #     # #   #
#     # #       #     # #       #  #
######  #       #     # #       ###
#     # #       #     # #       #  #
#     # #       #     # #     # #   #
######  ####### #######  #####  #    #
*/

func (r *Repository) CreateBlock(block Block) error {

	_, erri := r.db.Exec("INSERT INTO Blocks SET BlockBy = ?, BlockTo = ?", block.BlockBy, block.BlockTo)
	if erri != nil {
		return errors.New("Error adding block")
	}

	return nil
}

func (r *Repository) DeleteBlock(block Block) error {

	result, erri := r.db.Exec("DELETE FROM Blocks WHERE BlockBy = ?, BlockTo = ?", block.BlockBy, block.BlockTo)

	if erri != nil {
		return errors.New("Error deleting block")
	}

	rows, errr := result.RowsAffected()

	if errr != nil {
		return errors.New("Internal server error")
	}

	if rows == 0 {
		return errors.New("The block does not exist")
	}

	return nil
}

func (r *Repository) ExistsBlock(block Block) bool {
	var by string

	err := r.db.QueryRow("SELECT BlockBy FROM Users WHERE BlockBy = ? AND BlockTo = ?", block.BlockBy, block.BlockTo).Scan(&by)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	return true
}
