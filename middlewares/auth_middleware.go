package middlewares

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kampanosg/go-lsi/types"
)

type authMiddleware struct {
	signingKey []byte
}

var (
	logger = log.Default()

	errAuthFailed       = errors.New("auth failed")
	errUserUnauthorized = errors.New("user is not authorized")
	errBadDate          = errors.New("invalid date provided")
	errInvalidToken     = errors.New("token is invalid")
)

func NewAuthMiddleware(signKey []byte) *authMiddleware {
	return &authMiddleware{signingKey: signKey}
}

func (m *authMiddleware) ProtectedEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			logger.Printf("%s malformed auth header. header=%v\n", r.RequestURI, authHeader)
			unauthorised(w)
		} else {
			jwtToken := authHeader[1]
			token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				return m.signingKey, nil
			})

			if token != nil && token.Valid {
				next.ServeHTTP(w, r)
				return
			} else if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorMalformed != 0 {
					log.Printf("auth failed: that's not even a token. err=%v\n", err)
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					log.Printf("auth failed: token is either expired or not active yet. err=%v\n", err)
				} else {
					log.Printf("auth failed: couldn't handle this token. err=%v\n", err)
				}
			} else {
				log.Printf("auth failed: couldn't handle this token. err=%v\n", err)
			}
			unauthorised(w)
		}
	})
}

func unauthorised(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(types.ErrorResp{Message: errAuthFailed.Error(), Timestamp: time.Now()})
}
