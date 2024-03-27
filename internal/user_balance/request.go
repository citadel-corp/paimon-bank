package userbalance

import (
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

var imgUrlValidationRule = validation.NewStringRule(func(s string) bool {
	// match, _ := regexp.MatchString("^((http|https)://)[-a-zA-Z0-9@:%._\\+~#?&//=]{2,256}\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%._\\+~#?&//=]*)$", s)
	match, _ := regexp.MatchString(`^(http:\/\/www\.|https:\/\/www\.|http:\/\/|https:\/\/|\/|\/\/)?[A-z0-9_-]*?[:]?[A-z0-9_-]*?[@]?[A-z0-9]+([\-\.]{1}[a-z0-9]+)*\.[a-z]{2,5}(:[0-9]{1,5})?(\/{1}[A-z0-9_-]+)?(\/[A-z0-9_-]+)\.(?:jpg|jpeg|png)$`, s)
	return match
}, "image url is not valid")

type CreateUserBalancePayload struct {
	SenderBankAccountNumber string `json:"senderBankAccountNumber"`
	SenderBankName          string `json:"senderBankName"`
	AddedBalance            int    `json:"addedBalance"`
	Currency                string `json:"currency"`
	TransferProofImg        string `json:"transferProofImg"`
	UserID                  string
}

func (p CreateUserBalancePayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.SenderBankAccountNumber, validation.Required, validation.Length(5, 30)),
		validation.Field(&p.SenderBankName, validation.Required, validation.Length(5, 30)),
		validation.Field(&p.AddedBalance, validation.Required, validation.Min(0)),
		validation.Field(&p.Currency, validation.Required, is.CurrencyCode),
		validation.Field(&p.TransferProofImg, validation.Required, imgUrlValidationRule),
	)
}

type CreateTransactionPayload struct {
	RecipientBankAccountNumber string `json:"recipientBankAccountNumber"`
	RecipientBankName          string `json:"recipientBankName"`
	Balances                   int    `json:"balances"`
	FromCurrency               string `json:"fromCurrency"`
	UserID                     string
}

func (p CreateTransactionPayload) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.RecipientBankAccountNumber, validation.Required, validation.Length(5, 30)),
		validation.Field(&p.RecipientBankName, validation.Required, validation.Length(5, 30)),
		validation.Field(&p.Balances, validation.Required, validation.Min(0)),
		validation.Field(&p.FromCurrency, validation.Required, is.CurrencyCode),
	)
}

type ListUserBalancePayload struct {
	UserID string
}
