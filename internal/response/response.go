package response

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

const (
	StatusOK                  = 200
	StatusBadRequest          = 400
	StatusNotFound            = 404
	StatusInternalServerError = 500
)

var statusToMessage = map[int]string{
	StatusOK:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusInternalServerError: "Internal Server Error",
}

type Response struct {
	StatusCode int
	Headers    headers.Headers
	Trailers   headers.Headers
	Body       *bytes.Buffer

	isFlushed bool

	w io.Writer
}

func New(w io.Writer) *Response {
	return &Response{
		StatusCode: StatusOK,
		Headers:    headers.New(),
		Trailers:   headers.New(),
		Body:       bytes.NewBuffer(nil),

		w: w,
	}
}

func (r *Response) SetStatusCode(code int) {
	r.StatusCode = code
}

func (r *Response) Write(b []byte) {
	r.Body.Write(b)
}

func (r *Response) WriteChunked(b []byte) {
	_, _ = r.w.Write(b)
}

func (r *Response) WriteHeaders(w io.Writer) {
	for k, v := range r.Headers {
		_, _ = w.Write([]byte(k))
		_, _ = w.Write([]byte(": "))
		_, _ = w.Write([]byte(v))
		_, _ = w.Write([]byte("\r\n"))
	}
	_, _ = w.Write([]byte("\r\n"))
}

func (r *Response) Flush() {
	if r.isFlushed {
		return
	}

	_, _ = r.w.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.StatusCode, statusToMessage[r.StatusCode])))
	r.WriteHeaders(r.w)
	r.isFlushed = true
	if len(r.Body.Bytes()) > 0 {
		r.Headers.Set("Content-Length", strconv.Itoa(len(r.Body.Bytes())))
		_, _ = r.w.Write(r.Body.Bytes())
	}
}

func (r *Response) FlushTrailers() {
	for k, v := range r.Trailers {
		_, _ = r.w.Write([]byte(k))
		_, _ = r.w.Write([]byte(": "))
		_, _ = r.w.Write([]byte(v))
		_, _ = r.w.Write([]byte("\r\n"))
	}
	_, _ = r.w.Write([]byte("\r\n"))
}
