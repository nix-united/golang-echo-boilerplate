package handlers

import "echo-demo-project/server"

type PostHandler struct {
	server *server.Server
}

func NewPostHandler(server *server.Server) *PostHandler {
	return &PostHandler{server: server}
}
