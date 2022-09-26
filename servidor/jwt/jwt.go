package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
)

type CustomPayload struct {
	jwt.Payload
	ID int `json:"ID,omitempty"`
}

type Jwt struct {
	hs *jwt.HMACSHA
}

func (j *Jwt) Open() {
	j.hs = jwt.NewHS256([]byte(os.Getenv("jwtSecret")))
}

func (j *Jwt) CreateToken(ID int) (error, string) {
	now := time.Now()
	pl := CustomPayload{
		Payload: jwt.Payload{
			ExpirationTime: jwt.NumericDate(now.Add(24 * 30 * time.Hour)),
			NotBefore:      jwt.NumericDate(now),
			IssuedAt:       jwt.NumericDate(now),
		},
		ID: ID,
	}

	token, err := jwt.Sign(pl, j.hs)
	if err != nil {
		return errors.New("Error creating token"), ""
	}

	return nil, string(token)

}

func (j *Jwt) VerifyToken(token string) (error, CustomPayload) {
	var pl CustomPayload
	_, err := jwt.Verify([]byte(token), j.hs, &pl)
	if err != nil {
		return errors.New("Error verifying token"), CustomPayload{}
	}

	return nil, pl
}
