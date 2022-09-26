package domain

import (
	"github.com/soyReymundus/social/emailservice"
	"github.com/soyReymundus/social/imageservice"
	"github.com/soyReymundus/social/jwt"
	"github.com/soyReymundus/social/repository"
)

type Domain struct {
	persistence  Persistence
	emailService EmailService
	imageservice ImageService
	jwt          JwtInterface
	messages     chan C_message
}

type Persistence interface {
	Open() error
	CreateUser(user repository.User) error
	DeleteUser(ID int) (bool, error)
	GetUser(ID int) (repository.User, error)
	GetUserByCode(code string) (repository.User, error)
	GetUserByEmail(email string) (repository.User, error)
	UpdateUser(ID int, user repository.User) error
	ExistsUserByEmail(email string) bool
	ExistsUser(ID int) bool

	CreateBlock(block repository.Block) error
	DeleteBlock(block repository.Block) error
	ExistsBlock(block repository.Block) bool

	CreateChat(chat repository.Chat) error
	DeleteChat(ID int) (bool, error)
	GetChat(ID int) (repository.Chat, error)
	GetChatsByUser(UserID, limit int) ([]repository.Chat, error)
	GetChatByUsers(ID1, ID2 int) (repository.Chat, error)
	UpdateChat(ID int, chat repository.Chat) error

	CreatePost(post repository.Post) error
	DeletePost(ID int) (bool, error)
	GetPost(ID int) (repository.Post, error)
	GetPostByTitle(title string) (repository.Post, error)
	GetPosts(page int) ([]repository.Post, error)
	GetPostsByUser(UserID, limit int) ([]repository.Post, error)
	UpdatePost(ID int, post repository.Post) error

	CreateMessage(message repository.Message) error
	DeleteMessage(ID int) (bool, error)
	GetMessage(ID int) (repository.Message, error)
	GetMessageByTime(ChatID, date int) (repository.Message, error)
	GetMessagesByChat(ChatID, btw1, btw2 int) ([]repository.Message, error)
	UpdateMessage(ID int, message repository.Message) error
}

type EmailService interface {
	Open() error
	NoReply(To []string, Subject, Message string) error
}

type ImageService interface {
	Open() error
	Check(hash string) (bool, error)
	Delete(hash string) (bool, error)
	Create(img string) (string, error)
}

type JwtInterface interface {
	CreateToken(ID int) (error, string)
	VerifyToken(token string) (error, jwt.CustomPayload)
	Open()
}

type datas interface {
	[]repository.Chat | repository.Chat | repository.User | []repository.Post | repository.Post | []repository.Message | repository.Message | string
}

type Response[T datas] struct {
	Mesagge string `json:"message,omitempty"`
	Status  string `json:"status"`
	Data    T      `json:"data,omitempty"`
	Code    int    `json:"-,"`
}

type V_code struct {
	Code string `json:"code"`
}

type U_user struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	NewPassword string `json:"newPassword"`
	Username    string `json:"username"`
	Avatar      string `json:"avatar"`
	Status      int    `json:"status"`
}

type R_user struct {
	Password string `json:"password"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type L_user struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ID struct {
	Id int `json:"id"`
}

type T_post struct {
	Title string `json:"title"`
}

type R_post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type R_message struct {
	Content string `json:"content"`
}

type C_message struct {
	Message repository.Message `json:"message"`
	Action  string             `json:"action"`
}

type C_Error struct {
	Code    int    `json:"-,"`
	Message string `json:"message"`
}

func check(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func (d *Domain) Go() {
	d.persistence = &repository.Repository{}
	err := d.persistence.Open()
	check(err)

	d.emailService = &emailservice.EmailService{}
	err2 := d.persistence.Open()
	check(err2)

	d.imageservice = &imageservice.ImageService{}
	err3 := d.imageservice.Open()
	check(err3)

	d.jwt = &jwt.Jwt{}
	d.jwt.Open()

	d.messages = make(chan C_message)
}
