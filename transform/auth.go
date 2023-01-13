package transform

import (
	"kev/types/domain"
	"kev/types/response"
)

func FromAuthResponseToDomain(resp response.Auth) domain.Auth {
	return domain.Auth{
		Token:  resp.Token,
		UserId: resp.UserID,
		Server: resp.Server,
	}
}
