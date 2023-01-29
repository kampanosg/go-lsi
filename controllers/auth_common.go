package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/kampanosg/go-lsi/types"
)

var (
	errMethodNotSupported = errors.New("method not supported")
)

func failed(w http.ResponseWriter, err error, status int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(types.ErrorResp{Message: err.Error(), Timestamp: time.Now()})
}

func ok(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
