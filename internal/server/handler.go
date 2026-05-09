package server

import (
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
)

type Handler func(req *request.Request, resp *response.Response)
