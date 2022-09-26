package domain

import (
	"github.com/soyReymundus/social/repository"
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

func (d *Domain) OpenChat(id_user ID, token string) Response[repository.Chat] {
	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[repository.Chat]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	user, err := d.persistence.GetUser(payload.ID)

	if err != nil {
		return Response[repository.Chat]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if user.Banned {
		return Response[repository.Chat]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	_, errS := d.persistence.GetUser(id_user.Id)

	if errS != nil {
		if errS.Error() == "User does not exist" {
			return Response[repository.Chat]{Code: 404, Status: "error", Mesagge: "El usuario no existe"}
		} else {
			return Response[repository.Chat]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	chatG, errG := d.persistence.GetChatByUsers(payload.ID, id_user.Id)

	if errG == nil {
		return Response[repository.Chat]{Code: 200, Status: "success", Mesagge: "Conversacion abierta con exito", Data: chatG}
	} else {
		if errG.Error() != "The chat does not exist" {
			return Response[repository.Chat]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	if d.persistence.ExistsBlock(repository.Block{BlockTo: payload.ID, BlockBy: id_user.Id}) {
		return Response[repository.Chat]{Code: 403, Status: "error", Mesagge: "El otro usuario te ha bloqueado"}
	}
	if d.persistence.ExistsBlock(repository.Block{BlockTo: id_user.Id, BlockBy: payload.ID}) {
		return Response[repository.Chat]{Code: 403, Status: "error", Mesagge: "has bloqueado al otro usuario"}
	}

	if id_user.Id == payload.ID {
		return Response[repository.Chat]{Code: 409, Status: "error", Mesagge: "No se puede crear una conversacion con tu mismo"}
	}

	d.persistence.CreateChat(repository.Chat{
		UserID1: payload.ID,
		UserID2: id_user.Id,
		User1:   true,
		User2:   true,
	})

	//ajustar limite
	chat, errC := d.persistence.GetChatByUsers(payload.ID, id_user.Id)

	if errC != nil {
		return Response[repository.Chat]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[repository.Chat]{Code: 201, Status: "success", Mesagge: "Conversacion creada con exito", Data: chat}
}

func (d *Domain) CloseChat(id_chat ID, token string) Response[repository.Chat] {
	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[repository.Chat]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	user, err := d.persistence.GetUser(payload.ID)

	if err != nil {
		return Response[repository.Chat]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if user.Banned {
		return Response[repository.Chat]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	chat, errC := d.persistence.GetChat(id_chat.Id)

	if errC != nil {
		if errC.Error() == "The chat does not exist" {
			return Response[repository.Chat]{Code: 404, Status: "error", Mesagge: "La conversacion no existe"}
		} else {
			return Response[repository.Chat]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	if chat.UserID1 != payload.ID && chat.UserID2 != payload.ID {
		return Response[repository.Chat]{Code: 403, Status: "error", Mesagge: "No tienes autorizacion para hacer esto"}
	}

	if chat.UserID1 == payload.ID {
		chat.User1 = false
	}

	if chat.UserID2 == payload.ID {
		chat.User2 = false
	}

	errU := d.persistence.UpdateChat(id_chat.Id, chat)

	if errU != nil {
		return Response[repository.Chat]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[repository.Chat]{Code: 200, Status: "success", Mesagge: "Conversacion cerrada con exito", Data: chat}
}

func (d *Domain) GetChats(page int, token string) Response[[]repository.Chat] {
	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[[]repository.Chat]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	user, err := d.persistence.GetUser(payload.ID)

	if err != nil {
		return Response[[]repository.Chat]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if user.Banned {
		return Response[[]repository.Chat]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	chats, err := d.persistence.GetChatsByUser(payload.ID, page)

	if err != nil {
		return Response[[]repository.Chat]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[[]repository.Chat]{Code: 200, Status: "success", Mesagge: "Conversaciones obtenidas con exito", Data: chats}
}
