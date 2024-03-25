package user

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateUserPayload struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (p CreateUserPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Email, validation.Required, is.EmailFormat),
		validation.Field(&p.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 15)),
	)
}

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p LoginPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Email, validation.Required, is.EmailFormat),
		validation.Field(&p.Password, validation.Required, validation.Length(5, 15)),
	)
}
