package repository

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

/*
 #####  #     #    #    #######
#     # #     #   # #      #
#       #     #  #   #     #
#       ####### #     #    #
#       #     # #######    #
#     # #     # #     #    #
 #####  #     # #     #    #
*/

func (r *Repository) CreateChat(chat Chat) error {

	_, erri := r.db.Exec("INSERT INTO Chats SET UserID1 = ?, UserID2 = ?, User1 = ?, User2 = ?", chat.UserID1, chat.UserID2, chat.User1, chat.User2)
	if erri != nil {
		return errors.New("Error adding chat")
	}

	return nil
}

func (r *Repository) DeleteChat(ID int) (bool, error) {

	result, erri := r.db.Exec("DELETE FROM Chats WHERE ID = ?", ID)

	if erri != nil {
		return false, errors.New("Error deleting chat")
	}

	rows, errr := result.RowsAffected()

	if errr != nil {
		return false, errors.New("Internal server error")
	}

	if rows == 0 {
		return false, errors.New("The chat does not exist")
	}

	return true, nil
}

func (r *Repository) GetChat(ID int) (Chat, error) {

	var chat = Chat{}

	errs := r.db.QueryRow(`SELECT * FROM Chats WHERE ID = ?`, ID).Scan(&chat.UserID1, &chat.UserID2, &chat.User1, &chat.User2, &chat.ID)
	if errs == sql.ErrNoRows {
		return Chat{}, errors.New("The chat does not exist")
	} else if errs != nil {
		return Chat{}, errors.New("Error querying chat")
	}

	return chat, nil
}

func (r *Repository) GetChatByUsers(ID1, ID2 int) (Chat, error) {

	var chat = Chat{}

	errs := r.db.QueryRow(`SELECT * FROM Chats WHERE (UserID1 = ? OR UserID2 = ?) OR (UserID2 = ? OR UserID1 = ?)`, ID1, ID2, ID2, ID1).Scan(&chat.UserID1, &chat.UserID2, &chat.User1, &chat.User2, &chat.ID)
	if errs == sql.ErrNoRows {
		return Chat{}, errors.New("The chat does not exist")
	} else if errs != nil {
		return Chat{}, errors.New("Error querying chat")
	}

	return chat, nil
}

func (r *Repository) GetChatsByUser(UserID, page int) ([]Chat, error) {

	var chats = make([]Chat, 50)
	var max = 0
	r.db.QueryRow("SELECT MAX(ID) FROM Posts").Scan(&max)

	results, err := r.db.Query(`SELECT * FROM Chats WHERE (UserID1 = ? OR UserID2 = ?) AND ID BETWEEN ? AND ? ORDER BY ID DESC LIMIT = 50`, UserID, UserID, ((page-1)*50)-max, ((page)*50)-max)

	if err == sql.ErrNoRows {
		//Si no hay chats deberia estar todo bien
		return []Chat{}, nil
	} else if err != nil {
		return []Chat{}, errors.New("error querying chats")
	}

	for i := 0; results.Next(); i++ {
		chat := Chat{}

		err := results.Scan(&chat.UserID1, &chat.UserID2, &chat.User1, &chat.User2, &chat.ID)
		if err == sql.ErrNoRows {
			return []Chat{}, errors.New("Internal error")
		}

		chats[i] = chat
	}

	return chats, nil
}

func (r *Repository) UpdateChat(ID int, chat Chat) error {

	result, erru := r.db.Exec(`UPDATE Chats SET UserID1 = ?, UserID2 = ?, User1 = ?, User2 = ? WHERE ID = ?`, chat.UserID1, chat.UserID2, chat.User1, chat.User2, chat.ID)
	if erru != nil {
		return errors.New("Error updating chat")
	}

	rows, errr := result.RowsAffected()

	if errr != nil {
		return errors.New("Internal server error")
	}

	if rows == 0 {
		return errors.New("The chat does not exist")
	}

	return nil
}
