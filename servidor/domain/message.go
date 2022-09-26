package domain

import (
	"time"

	"github.com/soyReymundus/social/repository"
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

func (d *Domain) CreateMessage(message R_message, id_chat ID, token string) Response[repository.Message] {
	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[repository.Message]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	userT, errT := d.persistence.GetUser(payload.ID)

	if errT != nil {
		return Response[repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if userT.Banned {
		return Response[repository.Message]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	if len(message.Content) > 1000 {
		return Response[repository.Message]{Code: 422, Status: "error", Mesagge: "El mensaje debe ser menor a 1000 caracteres"}
	}

	chat, errC := d.persistence.GetChat(id_chat.Id)

	if errC != nil {
		if errC.Error() == "The chat does not exist" {
			return Response[repository.Message]{Code: 404, Status: "error", Mesagge: "La conversacion no existe"}
		} else {
			return Response[repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	if chat.UserID1 != payload.ID && chat.UserID2 != payload.ID {
		return Response[repository.Message]{Code: 403, Status: "error", Mesagge: "No tienes autorizacion para hacer esto"}
	}

	user := 0

	if chat.UserID1 == payload.ID {
		user = chat.UserID2
	} else {
		user = chat.UserID1
	}

	if d.persistence.ExistsBlock(repository.Block{BlockTo: payload.ID, BlockBy: user}) {
		return Response[repository.Message]{Code: 403, Status: "error", Mesagge: "El otro usuario te ha bloqueado"}
	}
	if d.persistence.ExistsBlock(repository.Block{BlockTo: user, BlockBy: payload.ID}) {
		return Response[repository.Message]{Code: 403, Status: "error", Mesagge: "has bloqueado al otro usuario"}
	}

	date := time.Now().Unix()

	errCC := d.persistence.CreateMessage(repository.Message{
		ChatID:    id_chat.Id,
		Timestamp: int(date),
		Content:   message.Content,
		AuthorID:  payload.ID,
	})

	if errCC != nil {
		return Response[repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	msg, errM := d.persistence.GetMessageByTime(id_chat.Id, int(date))

	if errM != nil {
		return Response[repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	d.messages <- C_message{
		Message: msg,
		Action:  "CREATE",
	}

	return Response[repository.Message]{Code: 201, Status: "success", Mesagge: "Mensaje creado con exito", Data: msg}
}

func (d *Domain) GetMessages(btw1, btw2 int, id_chat ID, token string) Response[[]repository.Message] {
	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[[]repository.Message]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	user, err := d.persistence.GetUser(payload.ID)

	if err != nil {
		return Response[[]repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if user.Banned {
		return Response[[]repository.Message]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	chat, errC := d.persistence.GetChat(id_chat.Id)

	if chat.UserID1 != payload.ID && chat.UserID2 != payload.ID {
		return Response[[]repository.Message]{Code: 403, Status: "error", Mesagge: "No tienes autorizacion para hacer esto"}
	}

	if errC != nil {
		if errC.Error() == "The chat does not exist" {
			return Response[[]repository.Message]{Code: 404, Status: "error", Mesagge: "La conversacion no existe"}
		} else {
			return Response[[]repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	if btw2-btw1 <= 0 {
		return Response[[]repository.Message]{Code: 406, Status: "error", Mesagge: "Debes pedir un numero positivo de mensajes"}
	}

	if btw2-btw1 > 100 {
		return Response[[]repository.Message]{Code: 422, Status: "error", Mesagge: "No puedes pedir mas de 100 mensajes"}
	}

	msgs, err := d.persistence.GetMessagesByChat(id_chat.Id, btw1, btw2)

	if err != nil {
		return Response[[]repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[[]repository.Message]{Code: 200, Status: "success", Mesagge: "Mensajes obtenidos con exito", Data: msgs}
}

func (d *Domain) GetMessage(id_message ID, token string) Response[repository.Message] {
	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[repository.Message]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	user, err := d.persistence.GetUser(payload.ID)

	if err != nil {
		return Response[repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if user.Banned {
		return Response[repository.Message]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	msg, errM := d.persistence.GetMessage(id_message.Id)

	if errM != nil {
		if errM.Error() == "The message does not exist" {
			return Response[repository.Message]{Code: 404, Status: "error", Mesagge: "El mensaje no existe"}
		} else {
			return Response[repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	chat, errC := d.persistence.GetChat(msg.ChatID)

	if errC != nil {
		return Response[repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if chat.UserID1 != payload.ID && chat.UserID2 != payload.ID {
		return Response[repository.Message]{Code: 403, Status: "error", Mesagge: "No tienes autorizacion para hacer esto"}
	}

	return Response[repository.Message]{Code: 200, Status: "success", Mesagge: "Mensaje obtenido con exito", Data: msg}
}

func (d *Domain) UpdateMessage(id_message ID, message R_message, token string) Response[repository.Message] {
	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[repository.Message]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	user, err := d.persistence.GetUser(payload.ID)

	if err != nil {
		return Response[repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if user.Banned {
		return Response[repository.Message]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	msg, errM := d.persistence.GetMessage(id_message.Id)

	if errM != nil {
		if errM.Error() == "The message does not exist" {
			return Response[repository.Message]{Code: 404, Status: "error", Mesagge: "El mensaje no existe"}
		} else {
			return Response[repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	if msg.AuthorID != payload.ID {
		return Response[repository.Message]{Code: 403, Status: "error", Mesagge: "No tienes autorizacion para hacer esto"}
	}

	if len(message.Content) > 1000 {
		return Response[repository.Message]{Code: 422, Status: "error", Mesagge: "El mensaje debe ser menor a 1000 caracteres"}
	}

	msg.Content = message.Content

	errU := d.persistence.UpdateMessage(msg.ID, msg)

	if errU != nil {
		return Response[repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	d.messages <- C_message{
		Message: msg,
		Action:  "UPDATE",
	}

	return Response[repository.Message]{Code: 200, Status: "success", Mesagge: "Mensaje obtenido con exito", Data: msg}
}

func (d *Domain) DeleteMessage(id_message ID, token string) Response[repository.Message] {
	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[repository.Message]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	user, err := d.persistence.GetUser(payload.ID)

	if err != nil {
		return Response[repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if user.Banned {
		return Response[repository.Message]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	msg, errM := d.persistence.GetMessage(id_message.Id)

	if errM != nil {
		if errM.Error() == "The message does not exist" {
			return Response[repository.Message]{Code: 404, Status: "error", Mesagge: "El mensaje no existe"}
		} else {
			return Response[repository.Message]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	if msg.AuthorID != payload.ID {
		return Response[repository.Message]{Code: 403, Status: "error", Mesagge: "No tienes autorizacion para hacer esto"}
	}

	d.messages <- C_message{
		Message: msg,
		Action:  "DELETE",
	}

	return Response[repository.Message]{Code: 200, Status: "success", Mesagge: "Mensaje eliminado con exito"}
}

func (d *Domain) GetMessageChannel(token string) (chan C_message, C_Error) {
	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return nil, C_Error{Code: 3000, Message: "Token no valido"}
	}

	user, errU := d.persistence.GetUser(payload.ID)

	if errU != nil {
		return nil, C_Error{Code: 1011, Message: "Error interno"}
	}

	if user.Banned {
		return nil, C_Error{Code: 1008, Message: "Tu cuenta esta bloqueada"}
	}

	user.StatusID = 1

	Uerr := d.persistence.UpdateUser(user.ID, user)

	if Uerr != nil {
		return nil, C_Error{Code: 1011, Message: "Error interno"}
	}

	c := make(chan C_message)

	go func() {

		defer func() {
			user.StatusID = 0
			d.persistence.UpdateUser(user.ID, user)

			recover()
		}()

		for {
			message, ok := <-d.messages

			if !ok {
				break
			}

			chat, err := d.persistence.GetChat(message.Message.ChatID)

			if err != nil {
				continue
			}

			if chat.UserID1 != payload.ID && chat.UserID2 != payload.ID {
				continue
			}

			c <- message
		}
	}()

	return c, C_Error{}
}
