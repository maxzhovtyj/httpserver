package server

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"net"
)

type Server struct {
	port int

	listener net.Listener

	handlers map[string]Handler
}

func New(port int) *Server {
	return &Server{
		port:     port,
		handlers: make(map[string]Handler),
	}
}

func (s *Server) Serve() error {
	var err error

	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	go s.listen()

	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Printf("Failed to close the connection%v\n", err)
			return
		}
	}()

	req, err := request.ParseFromReader(conn)
	if err != nil {
		resp := fmt.Sprintf("HTTP/1.1 400 BadRequest\r\nContent-Type: text/plain\r\n%s\r\n", err)
		_, err = conn.Write([]byte(resp))
		if err != nil {
			log.Printf("Failed to write to the connection%v\n", err)
		}

		return
	}

	resp := response.New(conn)
	resp.Headers = headers.NewDefault()

	handler, ok := s.handlers[req.Line.Target]
	if !ok {
		resp.StatusCode = response.StatusNotFound
		s.writeResponse(conn, resp)
		return
	}

	handler(req, resp)

	resp.Flush()
}

func (s *Server) writeResponse(conn net.Conn, resp *response.Response) {
	_, _ = conn.Write(resp.Body.Bytes())
}

func (s *Server) Close() error {
	return s.listener.Close()
}

func (s *Server) HandleFunc(target string, handler Handler) {
	s.handlers[target] = handler
}
