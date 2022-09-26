package domain

import (
	"github.com/soyReymundus/social/repository"
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

func (d *Domain) CreatePost(r_post R_post, token string) Response[repository.Post] {

	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[repository.Post]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	user, err := d.persistence.GetUser(payload.ID)

	if err != nil {
		return Response[repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if user.Banned {
		return Response[repository.Post]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	if r_post.Title == "" {
		return Response[repository.Post]{Code: 422, Status: "error", Mesagge: "Se necesita un titulo"}
	}

	if len(r_post.Title) > 30 {
		return Response[repository.Post]{Code: 422, Status: "error", Mesagge: "El titulo debe ser menor a 30 caracteres"}
	}

	_, errP := d.persistence.GetPostByTitle(r_post.Title)

	if errP == nil {
		return Response[repository.Post]{Code: 409, Status: "error", Mesagge: "El titulo ya esta en uso"}
	} else {
		if errP.Error() != "The post does not exist" {
			return Response[repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	if r_post.Content == "" {
		return Response[repository.Post]{Code: 422, Status: "error", Mesagge: "Se necesita contenido"}
	}

	if len(r_post.Title) > 10000 {
		return Response[repository.Post]{Code: 422, Status: "error", Mesagge: "El contenido debe ser menor a 10000 caracteres"}
	}

	errC := d.persistence.CreatePost(repository.Post{
		Title:      r_post.Title,
		Content:    r_post.Content,
		AuthorID:   payload.ID,
		HiddenPost: false,
	})

	if errC != nil {
		return Response[repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	post, errG := d.persistence.GetPostByTitle(r_post.Title)

	if errC != errG {
		return Response[repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[repository.Post]{Code: 201, Status: "success", Mesagge: "Publicacion creada con exito", Data: post}
}

func (d *Domain) GetPost(id_post ID) Response[repository.Post] {
	post, err := d.persistence.GetPost(id_post.Id)

	if err != nil {
		if err.Error() == "The post does not exist" {
			return Response[repository.Post]{Code: 404, Status: "error", Mesagge: "La publicacion no existe"}
		} else {
			return Response[repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	if post.HiddenPost {
		return Response[repository.Post]{Code: 423, Status: "error", Mesagge: "La publicacion fue eliminada"}
	}

	return Response[repository.Post]{Code: 200, Status: "success", Mesagge: "Publicacion obtenida", Data: post}
}

func (d *Domain) GetPosts(page int) Response[[]repository.Post] {
	posts, err := d.persistence.GetPosts(page)
	f := make([]repository.Post, len(posts))
	n := 0

	if err != nil {
		return Response[[]repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	for i, p := range posts {
		if !p.HiddenPost {
			f[n] = posts[i]
			n++
		}
	}

	return Response[[]repository.Post]{Code: 200, Status: "success", Mesagge: "Publicaciones obtenidas", Data: posts}
}

func (d *Domain) GetPostByTitle(t_post T_post) Response[repository.Post] {
	post, err := d.persistence.GetPostByTitle(t_post.Title)

	if err != nil {
		if err.Error() == "The post does not exist" {
			return Response[repository.Post]{Code: 404, Status: "error", Mesagge: "La publicacion no existe"}
		} else {
			return Response[repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	if post.HiddenPost {
		return Response[repository.Post]{Code: 423, Status: "error", Mesagge: "La publicacion fue eliminada"}
	}

	return Response[repository.Post]{Code: 200, Status: "success", Mesagge: "Publicacion obtenida", Data: post}
}

func (d *Domain) UpdatePost(id_post ID, r_post R_post, token string) Response[repository.Post] {
	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[repository.Post]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	post, errP := d.persistence.GetPost(id_post.Id)
	user, userErr := d.persistence.GetUser(payload.ID)

	if userErr != nil {
		return Response[repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if user.Banned {
		return Response[repository.Post]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	if errP != nil {
		if errP.Error() == "The post does not exist" {
			return Response[repository.Post]{Code: 404, Status: "error", Mesagge: "La publicacion no existe"}
		} else {
			return Response[repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	if post.HiddenPost {
		return Response[repository.Post]{Code: 423, Status: "error", Mesagge: "La publicacion fue eliminada"}
	}

	if user.ID != post.AuthorID {
		return Response[repository.Post]{Code: 403, Status: "error", Mesagge: "No tienes autorizacion para hacer esto"}
	}

	if r_post.Content == "" && r_post.Title == "" {
		return Response[repository.Post]{Code: 200, Status: "error", Mesagge: "No se efectuaron cambios", Data: post}
	}

	if r_post.Title != "" {
		if len(r_post.Title) > 30 {
			return Response[repository.Post]{Code: 422, Status: "error", Mesagge: "El titulo debe ser menor a 30 caracteres"}
		}

		_, errT := d.persistence.GetPostByTitle(r_post.Title)

		if errT == nil {
			return Response[repository.Post]{Code: 409, Status: "error", Mesagge: "El titulo ya esta en uso"}
		} else {
			if errT.Error() != "The post does not exist" {
				return Response[repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
			}
		}
		post.Title = r_post.Title
	}

	if r_post.Content != "" {
		if len(r_post.Title) > 10000 {
			return Response[repository.Post]{Code: 422, Status: "error", Mesagge: "El contenido debe ser menor a 10000 caracteres"}
		}

		post.Content = r_post.Content
	}

	errU := d.persistence.UpdatePost(post.ID, post)

	if errU != nil {
		return Response[repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[repository.Post]{Code: 200, Status: "success", Mesagge: "La publicacion se modifico con exito", Data: post}
}

func (d *Domain) HidePost(id_post ID, token string) Response[repository.Post] {
	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[repository.Post]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	post, errP := d.persistence.GetPost(id_post.Id)
	admin, admErr := d.persistence.GetUser(payload.ID)

	if admin.Banned {
		return Response[repository.Post]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	if admErr != nil {
		return Response[repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if errP != nil {
		if errP.Error() == "The post does not exist" {
			return Response[repository.Post]{Code: 404, Status: "error", Mesagge: "La publicacion no existe"}
		} else {
			return Response[repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	if !admin.Admin {
		return Response[repository.Post]{Code: 403, Status: "error", Mesagge: "No tienes autorizacion para hacer esto"}
	}

	post.HiddenPost = true

	return Response[repository.Post]{Code: 200, Status: "success", Mesagge: "La publicacion se oculto con exito", Data: post}
}

func (d *Domain) GetUserPosts(id_user ID) Response[[]repository.Post] {
	posts, err := d.persistence.GetPostsByUser(id_user.Id, 100)
	f := make([]repository.Post, len(posts))
	n := 0

	if err != nil {
		return Response[[]repository.Post]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	for i, p := range posts {
		if !p.HiddenPost {
			f[n] = posts[i]
			n++
		}
	}

	return Response[[]repository.Post]{Code: 200, Status: "success", Mesagge: "Publicaciones obtenidas", Data: posts}
}
