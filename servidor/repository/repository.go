package repository

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Estructura para guardar un usuario a la base de datos
type User struct {
	Username  string `json:"username"`
	Pass      string `json:"-,"`
	ID        int    `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	EmailCode string `json:"-,"`
	Avatar    string `json:"avatar,omitempty"`
	Admin     bool   `json:"admin"`
	Banned    bool   `json:"-,"`
	Verified  bool   `json:"-,"`
	StatusID  int8   `json:"status"`
}

// Estructura para guardar un estado a la base de datos
type StatusList struct {
	ID            int
	StatusMessage string
}

// Estructura para guardar un bloqueo a la base de datos
type Block struct {
	BlockTo int
	BlockBy int
}

// Estructura para guardar los chats a la base de datos
type Chat struct {
	UserID1 int  `json:"userID1,omitempty"`
	UserID2 int  `json:"userID2,omitempty"`
	User1   bool `json:"-,"`
	User2   bool `json:"-,"`
	ID      int  `json:"id,omitempty"`
}

// Estructura para guardar los mensajes a la base de datos
type Message struct {
	ChatID    int    `json:"chatID,omitempty"`
	AuthorID  int    `json:"authorID,omitempty"`
	Content   string `json:"content,omitempty"`
	Timestamp int    `json:"timestamp,omitempty"`
	ID        int    `json:"id,omitempty"`
}

// Estructura para guardar las publicaciones a la base de datos
type Post struct {
	Title      string `json:"title,omitempty"`
	AuthorID   int    `json:"authorID,omitempty"`
	Content    string `json:"content,omitempty"`
	HiddenPost bool   `json:"-,"`
	ID         int    `json:"id,omitempty"`
}

type Repository struct {
	db *sql.DB
}

func (r *Repository) Open() error {

	var err error

	r.db, err = sql.Open("mysql", os.Getenv("SQLUsername")+":"+os.Getenv("SQLPassword")+"@tcp("+os.Getenv("SQLHost")+":"+os.Getenv("SQLPort")+")/"+os.Getenv("SQLDB"))

	if err != nil {
		panic(err.Error())
	}

	return nil
}
