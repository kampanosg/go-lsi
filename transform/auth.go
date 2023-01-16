package transform

import (
	"github.com/kampanosg/go-lsi/types/domain"
	"github.com/kampanosg/go-lsi/types/response"
)

func FromAuthResponseToDomain(resp response.Auth) domain.Auth {
	return domain.Auth{
		Token:  resp.Token,
		UserId: resp.UserID,
		Server: resp.Server,
	}
}
