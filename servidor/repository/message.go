package repository

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

/*
#     # #######  #####   #####     #     #####  #######
##   ## #       #     # #     #   # #   #     # #
# # # # #       #       #        #   #  #       #
#  #  # #####    #####   #####  #     # #  #### #####
#     # #             #       # ####### #     # #
#     # #       #     # #     # #     # #     # #
#     # #######  #####   #####  #     #  #####  #######
*/

func (r *Repository) CreateMessage(message Message) error {

	_, erri := r.db.Exec("INSERT INTO Messages SET ChatID = ?,AuthorID = ?,Content = ?,Timestamp = ?", message.ChatID, message.AuthorID, message.Content, message.Timestamp)
	if erri != nil {
		return errors.New("Error adding message")
	}

	return nil
}

func (r *Repository) DeleteMessage(ID int) (bool, error) {

	result, erri := r.db.Exec("DELETE FROM Messages WHERE ID = ?", ID)

	if erri != nil {
		return false, errors.New("Error deleting message")
	}

	rows, errr := result.RowsAffected()

	if errr != nil {
		return false, errors.New("Internal server error")
	}

	if rows == 0 {
		return false, errors.New("The message does not exist")
	}

	return true, nil
}

func (r *Repository) GetMessage(ID int) (Message, error) {

	var message = Message{}

	errs := r.db.QueryRow(`SELECT * FROM Posts WHERE ID = ?`, ID).Scan(&message.ID, &message.ChatID, &message.AuthorID, &message.Content)
	if errs == sql.ErrNoRows {
		return Message{}, errors.New("The message does not exist")
	} else if errs != nil {
		return Message{}, errors.New("Error querying message")
	}

	return message, nil
}

func (r *Repository) GetMessagesByChat(ChatID, btw1, btw2 int) ([]Message, error) {

	var messages = make([]Message, btw2-btw1)

	results, err := r.db.Query(`SELECT * FROM Messages WHERE ChatID = ? AND ID BETWEEN ? AND ? ORDER BY ID DESC`, ChatID, btw1, btw2)

	if err == sql.ErrNoRows {
		//si no hay publicaciones deberia estar todo bien
		return []Message{}, nil
	} else if err != nil {
		return []Message{}, errors.New("error querying messages")
	}

	for i := 0; results.Next() && btw2-btw1 > i; i++ {
		message := Message{}

		err := results.Scan(&message.ID, &message.ChatID, &message.AuthorID, &message.Content)
		if err == sql.ErrNoRows {
			return []Message{}, errors.New("Internal error")
		}

		messages[i] = message
	}

	return messages, nil
}

func (r *Repository) GetMessageByTime(ChatID, date int) (Message, error) {
	var message = Message{}

	errs := r.db.QueryRow(`SELECT * FROM Posts WHERE ID = ? AND Timestamp = ?`, ChatID, date).Scan(&message.ID, &message.ChatID, &message.AuthorID, &message.Content)
	if errs == sql.ErrNoRows {
		return Message{}, errors.New("The message does not exist")
	} else if errs != nil {
		return Message{}, errors.New("Error querying message")
	}

	return message, nil
}

func (r *Repository) UpdateMessage(ID int, message Message) error {

	result, erru := r.db.Exec(`UPDATE Messages SET ChatID = ?,AuthorID = ?,Content = ?,Timestamp = ? WHERE ID = ?`, message.ChatID, message.AuthorID, message.Content, message.Timestamp, message.ID)
	if erru != nil {
		return errors.New("Error updating message")
	}

	rows, errr := result.RowsAffected()

	if errr != nil {
		return errors.New("Internal server error")
	}

	if rows == 0 {
		return errors.New("The message does not exist")
	}

	return nil
}
