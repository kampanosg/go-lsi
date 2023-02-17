package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kampanosg/go-lsi/clients/db"
	"github.com/kampanosg/go-lsi/types"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	errBadPassword = errors.New("username or password is incorrect")
)

const (
	TOKEN_LEN = 32
)

type AuthController struct {
	db         db.DB
	signingKey []byte
	logger     *zap.SugaredLogger
}

func NewAuthController(db db.DB, signKey []byte, logger *zap.SugaredLogger) AuthController {
	return AuthController{db: db, signingKey: signKey, logger: logger}
}

func (c *AuthController) HandleAuthRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		c.logger.Errorw("request failed", "reason", "method not supported", "method", r.Method, "uri", r.RequestURI)
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	var req types.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.logger.Errorw("request failed", "reason", "unable to decode body", "uri", r.RequestURI, "error", err.Error())
		failed(w, err, http.StatusBadRequest)
		return
	}

	user, err := c.db.GetUserByUsername(req.Username)
	if err != nil {
		c.logger.Errorw("request failed", "reason", "unable to auth user", "uri", r.RequestURI, "error", err.Error())
		failed(w, errBadPassword, http.StatusUnauthorized)
		return
	}

	if isInvalidPassword(req.Password, user.Password) {
		c.logger.Debugw("request failed", "reason", "bad username or password", "uri", r.RequestURI)
		failed(w, errBadPassword, http.StatusUnauthorized)
		return
	}

	now := time.Now()
	expiry := now.Add(time.Hour * 22 * 7)

	claims := &jwt.StandardClaims{
		ExpiresAt: expiry.Unix(),
		Subject:   user.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(c.signingKey)
	if err != nil {
		c.logger.Errorw("request failed", "reason", "unable to generate JWT", "user", req.Username, "uri", r.RequestURI, "error", err.Error())
		failed(w, errBadPassword, http.StatusUnauthorized)
		return
	}

	resp := types.AuthResp{
		Token:     tokenString,
		Message:   fmt.Sprintf("auth successful for user %s", req.Username),
		Timestamp: time.Now(),
	}

	c.logger.Infow("auth ok", "user", user.Username)

	ok(w, resp)
}

func (c *AuthController) HandlePasswordChangeRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		c.logger.Errorw("request failed", "reason", "method not supported", "method", r.Method, "uri", r.RequestURI)
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	var req types.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.logger.Errorw("request failed", "reason", "unable to decode body", "uri", r.RequestURI, "error", err.Error())
		failed(w, err, http.StatusBadRequest)
		return
	}

	user, err := c.db.GetUserByUsername(req.Username)
	if err != nil {
		c.logger.Errorw("request failed", "reason", "unable to auth user", "uri", r.RequestURI, "error", err.Error())
		failed(w, err, http.StatusUnauthorized)
		return
	}

	newPassword := strings.Trim(req.Password, " ")
	if len(newPassword) == 0 {
		c.logger.Errorw("request failed", "reason", "invalid password", "uri", r.RequestURI, "error", err.Error())
		failed(w, err, http.StatusBadRequest)
		return
	}

	hashedPword, err := hashPassword(newPassword)
	if err != nil {
		c.logger.Errorw("request failed", "reason", "unable to hash password", "uri", r.RequestURI, "error", err.Error())
		failed(w, err, http.StatusBadRequest)
		return
	}

	if err := c.db.UpdateUserPassword(user.ID, string(hashedPword)); err != nil {
		c.logger.Errorw("request failed", "reason", "unable to update password", "uri", r.RequestURI, "error", err.Error())
		failed(w, err, http.StatusBadRequest)
		return
	}

	resp := types.AuthResp{
		Token:     "login-again",
		Message:   fmt.Sprintf("password change successful for user %s", req.Username),
		Timestamp: time.Now(),
	}

	c.logger.Infow("password change ok", "user", user.Username)

	ok(w, resp)
}

func isInvalidPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err != nil
}

func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
