package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/kampanosg/go-lsi/clients/db"
	"github.com/kampanosg/go-lsi/types"
	"golang.org/x/crypto/bcrypt"
)

var (
	errBadPassword       = errors.New("username or password is incorrect")
)

const (
	TOKEN_LEN = 32
)

type AuthController struct {
	db         db.DB
	signingKey []byte
}

func NewAuthController(db db.DB, signKey []byte) AuthController {
	return AuthController{db: db, signingKey: signKey}
}

func (c *AuthController) HandleAuthRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		failed(w, errMethodNotSupported, http.StatusMethodNotAllowed)
		return
	}

	var req types.AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("%s - failed authenticating user. reason=%v\n", r.RequestURI, err)
		failed(w, err, http.StatusBadRequest)
		return
	}

	log.Printf("%s processing auth request, username=%s\n", r.RequestURI, req.Username)

	user, err := c.db.GetUserByUsername(req.Username)
	if err != nil {
		log.Printf("%s - unable to auth user=%s, reason=%v\n", r.RequestURI, req.Username, err)
		failed(w, errBadPassword, http.StatusUnauthorized)
		return
	}

	if isInvalidPassword(req.Password, user.Password) {
		log.Printf("%s - unable to auth user=%s, reason=%v\n", r.RequestURI, req.Username, errBadPassword)
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
		log.Printf("%s - unable to generate JWT for user=%s, reason=%v\n", r.RequestURI, req.Username, err)
		failed(w, errBadPassword, http.StatusUnauthorized)
		return
	}

	resp := types.AuthResponse{
		Token:     tokenString,
		Message:   fmt.Sprintf("auth successful for user %s", req.Username),
		Timestamp: time.Now(),
	}

	ok(w, resp)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func isInvalidPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err != nil
}
