package repository

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

/*
######  #######  #####  #######
#     # #     # #     #    #
#     # #     # #          #
######  #     #  #####     #
#       #     #       #    #
#       #     # #     #    #
#       #######  #####     #
*/
func (r *Repository) CreatePost(post Post) error {

	_, erri := r.db.Exec("INSERT INTO Posts SET AuthorID = ?,HiddenPost = ?,Title = ?,Content = ?", post.AuthorID, post.HiddenPost, post.Title, post.Content)
	if erri != nil {
		return errors.New("Error adding post")
	}

	return nil
}

func (r *Repository) DeletePost(ID int) (bool, error) {

	result, erri := r.db.Exec("DELETE FROM Posts WHERE ID = ?", ID)

	if erri != nil {
		return false, errors.New("Error deleting post")
	}

	rows, errr := result.RowsAffected()

	if errr != nil {
		return false, errors.New("Internal server error")
	}

	if rows == 0 {
		return false, errors.New("The post does not exist")
	}

	return true, nil
}

func (r *Repository) GetPost(ID int) (Post, error) {

	var post = Post{}

	errs := r.db.QueryRow(`SELECT * FROM Posts WHERE ID = ?`, ID).Scan(&post.ID, &post.AuthorID, &post.HiddenPost, &post.Title, &post.Content)
	if errs == sql.ErrNoRows {
		return Post{}, errors.New("The post does not exist")
	} else if errs != nil {
		return Post{}, errors.New("Error querying user")
	}

	return post, nil
}

func (r *Repository) GetPostByTitle(title string) (Post, error) {

	var post = Post{}

	errs := r.db.QueryRow(`SELECT * FROM Posts WHERE Title = ?`, title).Scan(&post.ID, &post.AuthorID, &post.HiddenPost, &post.Title, &post.Content)
	if errs == sql.ErrNoRows {
		return Post{}, errors.New("The post does not exist")
	} else if errs != nil {
		return Post{}, errors.New("Error querying user")
	}

	return post, nil
}

func (r *Repository) GetPosts(page int) ([]Post, error) {

	var posts = make([]Post, 50)
	var max = 0

	r.db.QueryRow("SELECT MAX(ID) FROM Posts").Scan(&max)

	results, err := r.db.Query(`SELECT * FROM Posts WHERE ID BETWEEN ? AND ? ORDER BY ID DESC LIMIT = 50`, ((page-1)*50)-max, ((page)*50)-max)

	if err == sql.ErrNoRows {
		//si no hay publicaciones deberia estar todo bien
		return []Post{}, nil
	} else if err != nil {
		return []Post{}, errors.New("error querying posts")
	}

	for i := 0; results.Next() && 50 > i; i++ {
		post := Post{}

		err := results.Scan(&post.ID, &post.AuthorID, &post.HiddenPost, &post.Title, &post.Content)
		if err == sql.ErrNoRows {
			return []Post{}, errors.New("Internal error")
		}

		posts[i] = post
	}

	return posts, nil
}

func (r *Repository) GetPostsByUser(UserID, limit int) ([]Post, error) {

	var posts = make([]Post, limit)

	results, err := r.db.Query(`SELECT * FROM Posts WHERE AuthorID = ? ORDER BY ID DESC LIMIT ?`, UserID, limit)

	if err == sql.ErrNoRows {
		//si no hay publicaciones deberia estar todo bien
		return []Post{}, nil
	} else if err != nil {
		return []Post{}, errors.New("error querying chats")
	}

	for i := 0; results.Next() && limit > i; i++ {
		post := Post{}

		err := results.Scan(&post.ID, &post.AuthorID, &post.HiddenPost, &post.Title, &post.Content)
		if err == sql.ErrNoRows {
			return []Post{}, errors.New("Internal error")
		}

		posts[i] = post
	}

	return posts, nil
}

func (r *Repository) UpdatePost(ID int, post Post) error {

	result, erru := r.db.Exec(`UPDATE Posts SET AuthorID = ?,HiddenPost = ?,Title = ?,Content = ? WHERE ID = ?`, post.AuthorID, post.HiddenPost, post.Title, post.Content, post.ID)
	if erru != nil {
		return errors.New("Error updating post")
	}

	rows, errr := result.RowsAffected()

	if errr != nil {
		return errors.New("Internal server error")
	}

	if rows == 0 {
		return errors.New("The post does not exist")
	}

	return nil
}
