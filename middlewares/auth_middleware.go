package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
)

type authMiddleware struct {
	signingKey []byte
	logger     *zap.SugaredLogger
}

var (
	errAuthFailed       = errors.New("auth failed")
	errUserUnauthorized = errors.New("user is not authorized")
	errBadDate          = errors.New("invalid date provided")
	errInvalidToken     = errors.New("token is invalid")
)

func NewAuthMiddleware(signKey []byte, logger *zap.SugaredLogger) *authMiddleware {
	return &authMiddleware{signingKey: signKey, logger: logger}
}

func (m *authMiddleware) ProtectedEndpoint(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			m.logger.Debugf("malformed auth header",
				"uri", r.RequestURI,
				"header", authHeader,
			)
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
					m.logger.Debugw("auth failed",
						"reason", "malformed token",
						"token", token,
						"error", err.Error(),
					)
				} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					m.logger.Debugw("auth failed",
						"reason", "token expired",
						"token", token,
						"error", err.Error(),
					)
				} else {
					m.logger.Debugw("auth failed",
						"reason", "couldn't handle this token",
						"token", token,
						"error", err.Error(),
					)
				}
			} else {
				m.logger.Debugw("auth failed",
					"reason", "couldn't handle this token",
					"token", token,
					"error", err.Error(),
				)
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
