package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"regexp"
	"strings"

	"github.com/soyReymundus/social/repository"
	"github.com/vincent-petithory/dataurl"
)

/*
#     #  #####  ####### ######
#     # #     # #       #     #
#     # #       #       #     #
#     #  #####  #####   ######
#     #       # #       #   #
#     # #     # #       #    #
 #####   #####  ####### #     #
*/

func (d *Domain) CreateUser(r_user R_user) Response[string] {

	if len(r_user.Username) > 20 {
		return Response[string]{Code: 422, Status: "error", Mesagge: "El nombre de usuario es muy largo"}
	}

	if len(r_user.Username) <= 3 {
		return Response[string]{Code: 422, Status: "error", Mesagge: "El nombre de usuario es muy corto"}
	}

	imgcheck, errCheck := d.imageservice.Check(r_user.Avatar)

	if errCheck != nil {
		return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if !imgcheck && r_user.Avatar != "" {
		return Response[string]{Code: 422, Status: "error", Mesagge: "La imagen no existe"}
	}

	match, _ := regexp.Match(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, []byte(r_user.Email))

	if len(r_user.Email) > 320 || !match {
		return Response[string]{Code: 422, Status: "error", Mesagge: "El correo electronico no es valido"}
	}

	if d.persistence.ExistsUserByEmail(strings.ToLower(r_user.Email)) {
		return Response[string]{Code: 409, Status: "error", Mesagge: "El correo electronico ya esta en uso"}
	}

	if r_user.Password == "" {
		return Response[string]{Code: 422, Status: "error", Mesagge: "Ingresa una contraseña"}
	}

	pass := sha256.Sum256([]byte(os.Getenv("salt") + r_user.Password))
	code := sha256.Sum256([]byte(os.Getenv("salt") + r_user.Email))

	user := repository.User{
		Username:  r_user.Username,
		Pass:      hex.EncodeToString(pass[:]),
		Email:     strings.ToLower(r_user.Email),
		EmailCode: hex.EncodeToString(code[:]),
		Verified:  false,
		Admin:     false,
		Avatar:    r_user.Avatar,
		Banned:    false,
		StatusID:  1,
	}

	err := d.persistence.CreateUser(user)

	if err != nil {
		return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	var emails = make([]string, 1)
	emails[0] = r_user.Email

	Eerr := d.emailService.NoReply(emails, "Su codigo de verificacion", "Entre aqui para verificar su cuenta: "+os.Getenv("WEB")+"/verify?code="+user.EmailCode)

	if Eerr != nil {
		return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[string]{Code: 201, Status: "success", Mesagge: "Usuario creador con exito, se envio un codigo de verificacion"}
}

func (d *Domain) CreateUserAvatar(data string) Response[string] {
	dataURL, err := dataurl.DecodeString(data)
	if err != nil {
		return Response[string]{Code: 400, Status: "error", Mesagge: "Los datos enviados no son validos"}
	}

	if dataURL.ContentType() != "image/gif" && dataURL.ContentType() != "image/webp" && dataURL.ContentType() != "image/png" && dataURL.ContentType() != "image/jpeg" && dataURL.ContentType() != "image/vnd.microsoft.icon" {
		return Response[string]{Code: 415, Status: "error", Mesagge: "Los datos enviandos no son una imagen valida"}
	}

	hash, errC := d.imageservice.Create(string(dataURL.Data))

	if errC != nil {
		return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[string]{Code: 201, Status: "success", Mesagge: "Avatar creado con exito", Data: hash}
}

func (d *Domain) VerifyUser(v_code V_code) Response[string] {

	if len(v_code.Code) != 64 {
		return Response[string]{Code: 422, Status: "error", Mesagge: "El codigo no es valido"}
	}

	user, err := d.persistence.GetUserByCode(v_code.Code)

	if err != nil {
		if err.Error() == "User does not exist" {
			return Response[string]{Code: 404, Status: "error", Mesagge: "El usuario no existe"}
		} else {
			return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	user.Verified = true
	user.EmailCode = ""

	erru := d.persistence.UpdateUser(user.ID, user)

	if erru != nil {
		return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[string]{Code: 200, Status: "success", Mesagge: "Usuario verificado con exito"}
}

func (d *Domain) UpdateUser(u_user U_user, token string) Response[repository.User] {
	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[repository.User]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	if payload.ID != u_user.Id {
		return Response[repository.User]{Code: 403, Status: "error", Mesagge: "No tienes autorizacion para hacer esto"}
	}

	user, err := d.persistence.GetUser(payload.ID)

	if user.Banned {
		return Response[repository.User]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	pass := sha256.Sum256([]byte(os.Getenv("salt") + u_user.Password))

	if err != nil {
		return Response[repository.User]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if hex.EncodeToString(pass[:]) != user.Pass {
		return Response[repository.User]{Code: 401, Status: "error", Mesagge: "Contraseña incorrecta"}
	}

	if u_user.NewPassword != "" {
		Newpass := sha256.Sum256([]byte(os.Getenv("salt") + u_user.NewPassword))
		user.Pass = string(Newpass[:])
	}

	if u_user.Username != "" {

		if len(u_user.Username) > 20 {
			return Response[repository.User]{Code: 422, Status: "error", Mesagge: "El nombre de usuario es muy largo"}
		}

		if len(u_user.Username) <= 3 {
			return Response[repository.User]{Code: 422, Status: "error", Mesagge: "El nombre de usuario es muy corto"}
		}

		user.Username = u_user.Username
	}

	if u_user.Avatar != "" {

		imgcheck, errCheck := d.imageservice.Check(u_user.Avatar)

		if errCheck != nil {
			return Response[repository.User]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}

		if !imgcheck {
			return Response[repository.User]{Code: 422, Status: "error", Mesagge: "La imagen no existe"}
		}

		user.Avatar = u_user.Avatar
	}

	if u_user.Email != "" {

		match, _ := regexp.Match(`^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, []byte(u_user.Email))

		if len(u_user.Email) > 320 || !match {
			return Response[repository.User]{Code: 422, Status: "error", Mesagge: "El correo electronico no es valido"}
		}

		if d.persistence.ExistsUserByEmail(strings.ToLower(u_user.Email)) {
			return Response[repository.User]{Code: 409, Status: "error", Mesagge: "El correo electronico ya esta en uso"}
		}

		user.Email = u_user.Email
		code := sha256.Sum256([]byte(os.Getenv("salt") + u_user.Email))
		user.EmailCode = string(code[:])
		user.Verified = false
		var emails = make([]string, 1)
		emails[0] = u_user.Email

		Eerr := d.emailService.NoReply(emails, "Su codigo de verificacion", "Entre aqui para verificar su cuenta: "+os.Getenv("WEB")+"/verify?code="+user.EmailCode)

		if Eerr != nil {
			return Response[repository.User]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	if u_user.Status > 2 || u_user.Status < 0 {
		return Response[repository.User]{Code: 422, Status: "error", Mesagge: "Ingresa un estado valido"}
	} else {
		user.StatusID = int8(u_user.Status)
	}

	Uerr := d.persistence.UpdateUser(user.ID, user)

	if Uerr != nil {
		return Response[repository.User]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[repository.User]{Code: 200, Status: "success", Mesagge: "El usuario se actualizo con exito", Data: user}
}

func (d *Domain) GetUser(id_user ID) Response[repository.User] {

	user, err := d.persistence.GetUser(id_user.Id)

	if err != nil {
		if err.Error() == "User does not exist" {
			return Response[repository.User]{Code: 404, Status: "error", Mesagge: "El usuario no existe"}
		} else {
			return Response[repository.User]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	return Response[repository.User]{Code: 200, Status: "success", Mesagge: "Se obtuvo el usuario con exito", Data: user}
}

func (d *Domain) GetMe(token string) Response[repository.User] {

	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[repository.User]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	user, err := d.persistence.GetUser(payload.ID)

	if err != nil {
		return Response[repository.User]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[repository.User]{Code: 200, Status: "success", Mesagge: "Se obtuvo tu usuario con exito", Data: user}
}

func (d *Domain) BanUser(id_user ID, token string) Response[string] {

	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[string]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	admin, admErr := d.persistence.GetUser(payload.ID)

	if admErr != nil {
		return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if admin.Banned {
		return Response[string]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	if !admin.Admin {
		return Response[string]{Code: 403, Status: "error", Mesagge: "No tienes autorizacion para hacer esto"}
	}

	user, err := d.persistence.GetUser(id_user.Id)

	if err != nil {
		if err.Error() == "User does not exist" {
			return Response[string]{Code: 404, Status: "error", Mesagge: "El usuario no existe"}
		} else {
			return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	user.Banned = true

	Uerr := d.persistence.UpdateUser(user.ID, user)

	if Uerr != nil {
		return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[string]{Code: 200, Status: "success", Mesagge: "El usuario se proibio con exito"}
}

func (d *Domain) CreateToken(l_user L_user) Response[string] {

	user, err := d.persistence.GetUserByEmail(l_user.Email)

	if err != nil {
		if err.Error() == "User does not exist" {
			return Response[string]{Code: 404, Status: "error", Mesagge: "El usuario no existe"}
		} else {
			return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	pass := sha256.Sum256([]byte(os.Getenv("salt") + l_user.Password))

	if user.Pass != hex.EncodeToString(pass[:]) {
		return Response[string]{Code: 401, Status: "error", Mesagge: "Contraseña incorrecta"}
	}

	if user.Banned {
		return Response[string]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	jwterr, token := d.jwt.CreateToken(user.ID)

	if jwterr != nil {
		return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[string]{Code: 200, Status: "success", Mesagge: "Sesion iniciada con exito", Data: token}
}

func (d *Domain) BlockUser(id_user ID, token string) Response[string] {

	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[string]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	userT, errT := d.persistence.GetUser(payload.ID)

	if errT != nil {
		return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if userT.Banned {
		return Response[string]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	user, err := d.persistence.GetUser(id_user.Id)

	if err != nil {
		if err.Error() == "User does not exist" {
			return Response[string]{Code: 404, Status: "error", Mesagge: "El usuario no existe"}
		} else {
			return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	block := repository.Block{
		BlockTo: user.ID,
		BlockBy: payload.ID,
	}

	if d.persistence.ExistsBlock(block) {
		return Response[string]{Code: 200, Status: "error", Mesagge: "El usuario ya fue bloqueado"}
	}

	errB := d.persistence.CreateBlock(block)

	if errB != nil {
		return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[string]{Code: 200, Status: "success", Mesagge: "El usuario fue bloqueado con exito"}
}

func (d *Domain) UnblockUser(id_user ID, token string) Response[string] {

	tokenerr, payload := d.jwt.VerifyToken(token)

	if tokenerr != nil {
		return Response[string]{Code: 400, Status: "error", Mesagge: "Token no valido"}
	}

	userT, errT := d.persistence.GetUser(payload.ID)

	if errT != nil {
		return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	if userT.Banned {
		return Response[string]{Code: 403, Status: "error", Mesagge: "Tu cuenta esta bloqueada"}
	}

	user, err := d.persistence.GetUser(id_user.Id)

	if err != nil {
		if err.Error() == "User does not exist" {
			return Response[string]{Code: 404, Status: "error", Mesagge: "El usuario no existe"}
		} else {
			return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
		}
	}

	block := repository.Block{
		BlockTo: user.ID,
		BlockBy: payload.ID,
	}

	if !d.persistence.ExistsBlock(block) {
		return Response[string]{Code: 200, Status: "error", Mesagge: "El usuario no esta bloqueado"}
	}

	errB := d.persistence.DeleteBlock(block)

	if errB != nil {
		return Response[string]{Code: 500, Status: "error", Mesagge: "Error interno"}
	}

	return Response[string]{Code: 200, Status: "success", Mesagge: "El usuario fue desbloqueado con exito"}
}
