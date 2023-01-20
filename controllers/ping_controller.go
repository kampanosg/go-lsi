package controllers

import (
	"net/http"
)

type PingController struct {}

func NewPingController() PingController {
	return PingController{}
}

func (c *PingController) HandlePingRequest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("PONG"))
}