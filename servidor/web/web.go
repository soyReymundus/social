package web

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/soyReymundus/social/domain"
)

type Server struct {
	app    *fiber.App
	domain domain.Domain
}

func (s *Server) createErrorResponse(code int, msg string) domain.Response[string] {
	return domain.Response[string]{
		Mesagge: msg,
		Status:  "error",
		Code:    code,
	}
}

func (s *Server) Open(dmn *domain.Domain) {
	s.domain = *dmn

	s.app = fiber.New(fiber.Config{

		BodyLimit: 5120,

		ErrorHandler: func(c *fiber.Ctx, err error) error {

			c.Type("json")

			if err.Error() == fiber.ErrRequestEntityTooLarge.Error() {
				errRes := s.createErrorResponse(413, "Cuerpo demasiando largo")
				c.Status(errRes.Code).JSON(errRes)
				return nil
			}

			errRes := s.createErrorResponse(500, "error interno")
			c.Status(errRes.Code).JSON(errRes)

			return nil
		},
	})

	s.app.Use(func(c *fiber.Ctx) error {

		c.Type("json")

		//verificar el Authorization
		if c.Get("Authorization") != "" {
			Authorization := strings.Split(c.Get("Authorization"), " ")

			if len(Authorization) != 2 {
				errRes := s.createErrorResponse(400, "La cabezera Authorization esta malformada")
				return c.Status(errRes.Code).JSON(errRes)
			}

			if len(strings.Split(Authorization[1], ".")) != 3 {
				errRes := s.createErrorResponse(400, "El token esta malformado")
				return c.Status(errRes.Code).JSON(errRes)
			}

			if Authorization[0] != "Bearer" {
				errRes := s.createErrorResponse(400, "El tipo de token enviado no esta aceptado")
				return c.Status(errRes.Code).JSON(errRes)
			}

			c.Locals("token", Authorization[1])
		} else {
			c.Locals("token", "")
		}

		if c.Method() != "GET" && c.Method() != "DELETE" && c.Method() != "HEAD" {

			//verificar el Content-Type
			if strings.ToLower(c.Get("Content-Type")) == "" {
				errRes := s.createErrorResponse(400, "Content-Type no definido")
				return c.Status(errRes.Code).JSON(errRes)
			}

			if strings.ToLower(c.Get("Content-Type")) != "application/json" {
				errRes := s.createErrorResponse(415, "Content-Type no aceptado")
				return c.Status(errRes.Code).JSON(errRes)
			}

			//verificar el Content-Length
			if c.Get("Content-Length") == "" {
				errRes := s.createErrorResponse(411, "Content-Length no definido")
				return c.Status(errRes.Code).JSON(errRes)
			}

			cl, err := strconv.Atoi(c.Get("Content-Length"))

			if err != nil {
				errRes := s.createErrorResponse(400, "Content-Length malformado")
				return c.Status(errRes.Code).JSON(errRes)
			}

			if len(c.Body()) != cl {
				errRes := s.createErrorResponse(400, "El Content-Length y el Body no coinciden")
				return c.Status(errRes.Code).JSON(errRes)
			}
		}

		return c.Next()
	})

	s.app.Route("/users", func(router fiber.Router) {
		router.Post("/", func(c *fiber.Ctx) error {

			u := new(domain.R_user)
			if err := c.BodyParser(u); err != nil {
				return err
			}

			r := s.domain.CreateUser(*u)

			return c.Status(r.Code).JSON(r)
		})

		router.Post("/auth", func(c *fiber.Ctx) error {

			u := new(domain.L_user)
			if err := c.BodyParser(u); err != nil {
				return err
			}

			r := s.domain.CreateToken(*u)

			return c.Status(r.Code).JSON(r)
		})

		router.Post("/pfp", func(c *fiber.Ctx) error {
			p := s.domain.CreateUserAvatar(string(c.Body()))
			return c.Status(p.Code).JSON(p)
		})

		router.Patch("/verify", func(c *fiber.Ctx) error {

			codeRes := new(domain.V_code)
			if err := c.BodyParser(codeRes); err != nil {
				return err
			}

			r := s.domain.VerifyUser(*codeRes)

			return c.Status(r.Code).JSON(r)
		})

		router.Get("/me", func(c *fiber.Ctx) error {

			u := s.domain.GetMe(c.Locals("token").(string))

			return c.Status(u.Code).JSON(u)
		})

		router.Get("/:ID", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Params("ID"))

			if err != nil {
				errRes := s.createErrorResponse(400, "La ID enviada no es valida")
				return c.Status(errRes.Code).JSON(errRes)
			}

			r := s.domain.GetUser(domain.ID{Id: id})

			return c.Status(r.Code).JSON(r)
		})

		router.Patch("/:ID", func(c *fiber.Ctx) error {

			u := new(domain.U_user)
			if err := c.BodyParser(u); err != nil {
				return err
			}

			r := s.domain.UpdateUser(*u, c.Locals("token").(string))

			return c.Status(r.Code).JSON(r)
		})

		router.Delete("/:ID", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Params("ID"))

			if err != nil {
				errRes := s.createErrorResponse(400, "La ID enviada no es valida")
				return c.Status(errRes.Code).JSON(errRes)
			}

			r := s.domain.BanUser(domain.ID{Id: id}, c.Locals("token").(string))

			return c.Status(r.Code).JSON(r)
		})

		router.Get("/:ID/post", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Params("ID"))

			if err != nil {
				errRes := s.createErrorResponse(400, "La ID enviada no es valida")
				return c.Status(errRes.Code).JSON(errRes)
			}

			r := s.domain.GetUserPosts(domain.ID{Id: id})

			return c.Status(r.Code).JSON(r)
		})
	})

	s.app.Route("/posts", func(router fiber.Router) {
		router.Get("/", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Query("page"))

			if err != nil {
				errRes := s.createErrorResponse(400, "La pagina enviada no es valida")
				return c.Status(errRes.Code).JSON(errRes)
			}

			r := s.domain.GetPosts(id)

			return c.Status(r.Code).JSON(r)
		})

		router.Post("/", func(c *fiber.Ctx) error {

			p := new(domain.R_post)
			if err := c.BodyParser(p); err != nil {
				return err
			}

			r := s.domain.CreatePost(*p, c.Locals("token").(string))

			return c.Status(r.Code).JSON(r)
		})

		router.Get("/:ID", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Params("ID"))

			if err != nil {
				r := s.domain.GetPostByTitle(domain.T_post{Title: c.Params("id")})

				return c.Status(r.Code).JSON(r)
			} else {
				r := s.domain.GetPost(domain.ID{Id: id})

				return c.Status(r.Code).JSON(r)
			}
		})

		router.Patch("/:ID", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Params("ID"))

			if err != nil {
				id = s.domain.GetPostByTitle(domain.T_post{Title: c.Params("id")}).Data.ID
			}

			p := new(domain.R_post)
			if err := c.BodyParser(p); err != nil {
				return err
			}

			r := s.domain.UpdatePost(domain.ID{Id: id}, *p, c.Locals("token").(string))

			return c.Status(r.Code).JSON(r)
		})

		router.Patch("/:ID", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Params("ID"))

			if err != nil {
				id = s.domain.GetPostByTitle(domain.T_post{Title: c.Params("id")}).Data.ID
			}

			r := s.domain.HidePost(domain.ID{Id: id}, c.Locals("token").(string))

			return c.Status(r.Code).JSON(r)
		})

		router.Delete("/:ID", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Params("ID"))

			if err != nil {
				id = s.domain.GetPostByTitle(domain.T_post{Title: c.Params("id")}).Data.ID
			}

			r := s.domain.HidePost(domain.ID{Id: id}, c.Locals("token").(string))

			return c.Status(r.Code).JSON(r)
		})

	})

	s.app.Route("/chats", func(router fiber.Router) {
		router.Get("/", func(c *fiber.Ctx) error {

			fmt.Println(websocket.IsWebSocketUpgrade(c))
			if websocket.IsWebSocketUpgrade(c) {
				return c.Next()
			}

			page, err := strconv.Atoi(c.Query("page"))

			if err != nil {
				errRes := s.createErrorResponse(400, "La pagina enviada no es valida")
				return c.Status(errRes.Code).JSON(errRes)

			}

			r := s.domain.GetChats(page, c.Locals("token").(string))

			return c.Status(r.Code).JSON(r)
		})

		router.Get("/", websocket.New(func(c *websocket.Conn) {
			messages, errG := s.domain.GetMessageChannel(c.Locals("token").(string))

			if errG.Code != 0 {
				out, _ := json.Marshal(errG)
				c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(errG.Code, string(out)))
				c.Close()
			}

			var (
				err     error
				message domain.C_message
				ok      bool
			)

			c.SetCloseHandler(func(code int, text string) error {
				close(messages)
				return nil
			})

			for {
				message, ok = <-messages

				if !ok {
					break
				}

				if err = c.WriteJSON(message); err != nil {
					out, _ := json.Marshal(domain.C_Error{Message: "Error interno"})
					c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, string(out)))
					c.Close()
					break
				}
			}
		}))

		router.Post("/", func(c *fiber.Ctx) error {

			u := new(domain.ID)
			if err := c.BodyParser(u); err != nil {
				return err
			}

			r := s.domain.OpenChat(*u, c.Locals("token").(string))

			return c.Status(r.Code).JSON(r)
		})

		router.Get("/:ID", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Params("ID"))

			btw1, errB1 := strconv.Atoi(c.Query("btw1"))
			btw2, errB2 := strconv.Atoi(c.Query("btw2"))

			if err != nil {
				errRes := s.createErrorResponse(400, "La ID enviada no es valida")
				return c.Status(errRes.Code).JSON(errRes)
			}

			if errB2 != nil || errB1 != nil {
				errRes := s.createErrorResponse(400, "El parametro enviando no es valido")
				return c.Status(errRes.Code).JSON(errRes)
			}

			r := s.domain.GetMessages(btw1, btw2, domain.ID{Id: id}, c.Locals("token").(string))

			return c.Status(r.Code).JSON(r)
		})

		router.Post("/:ID", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Params("ID"))

			if err != nil {
				errRes := s.createErrorResponse(400, "La ID enviada no es valida")
				return c.Status(errRes.Code).JSON(errRes)
			}

			m := new(domain.R_message)
			if err := c.BodyParser(m); err != nil {
				return err
			}

			r := s.domain.CreateMessage(*m, domain.ID{Id: id}, c.Locals("token").(string))

			return c.Status(r.Code).JSON(r)
		})

		router.Delete("/:ID", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Params("ID"))

			if err != nil {
				errRes := s.createErrorResponse(400, "La ID enviada no es valida")
				return c.Status(errRes.Code).JSON(errRes)
			}

			r := s.domain.CloseChat(domain.ID{Id: id}, c.Locals("token").(string))

			return c.Status(r.Code).JSON(r)
		})

		router.Get("/:ID/:msgID", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Params("ID"))
			msgID, errmsg := strconv.Atoi(c.Params("ID"))

			if err != nil {
				errRes := s.createErrorResponse(404, "La ID enviada no es valida")
				return c.Status(errRes.Code).JSON(errRes)
			}

			if errmsg != nil {
				errRes := s.createErrorResponse(404, "La ID del mensaje enviada no es valida")
				return c.Status(errRes.Code).JSON(errRes)
			}

			r := s.domain.GetMessage(domain.ID{Id: msgID}, c.Locals("token").(string))

			if r.Data.ChatID != id {
				errRes := s.createErrorResponse(404, "No se encontro el mensaje")
				return c.Status(errRes.Code).JSON(errRes)
			}

			return c.Status(r.Code).JSON(r)
		})

		router.Patch("/:ID/:msgID", func(c *fiber.Ctx) error {

			id, err := strconv.Atoi(c.Params("ID"))
			msgID, errmsg := strconv.Atoi(c.Params("ID"))

			if err != nil {
				errRes := s.createErrorResponse(404, "La ID enviada no es valida")
				return c.Status(errRes.Code).JSON(errRes)
			}

			if errmsg != nil {
				errRes := s.createErrorResponse(404, "La ID del mensaje enviada no es valida")
				return c.Status(errRes.Code).JSON(errRes)
			}

			m := new(domain.R_message)
			if err := c.BodyParser(m); err != nil {
				return err
			}

			check := s.domain.GetMessage(domain.ID{Id: msgID}, c.Locals("token").(string))

			if check.Data.ChatID != id {
				errRes := s.createErrorResponse(404, "No se encontro el mensaje")
				return c.Status(errRes.Code).JSON(errRes)
			}

			r := s.domain.UpdateMessage(domain.ID{Id: msgID}, *m, c.Locals("token").(string))

			return c.Status(r.Code).JSON(r)
		})

	})

	s.app.All("*", func(c *fiber.Ctx) error {
		errRes := s.createErrorResponse(404, "El endpoint no existe")
		return c.Status(errRes.Code).JSON(errRes)
	})

	s.app.Listen(os.Getenv("HTTPHost"))
}
