package request

import (
	"bytes"
	"errors"
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
	"strings"
)

type state int

const (
	initializedState state = iota
	requestStateParsingHeaders
	parsinBody
	doneState
)

type Request struct {
	Line    Line
	Headers headers.Headers
	Body    []byte

	state state
	buf   []byte

	currLine       []byte
	rawRequestLine []byte

	rawHeaders []byte

	bytesRead int
}

func (r *Request) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf(`%s %s HTTP/%s`, r.Line.Method, r.Line.Target, r.Line.Version))
	b.WriteString("\n")

	for k, v := range r.Headers {
		b.WriteString(k)
		b.WriteString(": ")
		b.WriteString(v)
		b.WriteString("\n")
	}

	b.Write(r.Body)

	return b.String()
}

type Line struct {
	Method  string
	Target  string
	Version string
}

func ParseFromReader(r io.Reader) (*Request, error) {
	req := new(Request)

	req.Headers = headers.New()
	req.currLine = make([]byte, 8)
	buf := make([]byte, 8)

	for {
		if req.state == initializedState {
			n, err := req.parseRequestLine(r)
			req.bytesRead += n
			if err != nil {
				return nil, err
			}

			continue
		}

		if req.state == requestStateParsingHeaders {
			if len(req.buf) > 0 {
				req.rawHeaders = append(req.rawHeaders, req.buf...)
				req.buf = req.buf[:0]
			}

			n, err := req.parseHeaders(r)
			req.bytesRead += n
			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}

				return nil, err
			}

			if req.state == parsinBody {
				contentLenRaw := req.Headers["content-length"]
				if contentLenRaw == "" {
					return req, nil
				}

				d, _ := strconv.Atoi(contentLenRaw)
				if d == 0 {
					return nil, nil
				}
			}

			continue
		}

		n, err := r.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		}

		req.Body = append(req.Body, buf[:n]...)
	}

	contentLenRaw := req.Headers["content-length"]
	d, _ := strconv.Atoi(contentLenRaw)
	if d == 0 {
		return req, nil
	}

	if d != len(req.Body) {
		return nil, fmt.Errorf("invalid body length")
	}

	return req, nil
}

// ParseLine parses a request line into a Line struct
// GET /coffee HTTP/1.1
//
// Method GET
// Target /coffee
// Version HTTP/1.1
func (r *Request) parseRequestLine(reader io.Reader) (int, error) {
	buf := r.currLine

	n, err := reader.Read(buf)
	if err != nil {
		return 0, err
	}

	if n == 0 {
		return 0, nil
	}

	buf = buf[:n]

	r.rawRequestLine = append(r.rawRequestLine, buf...)

	idx := bytes.Index(r.rawRequestLine, []byte("\r\n"))
	if idx != -1 {
		r.buf = append(r.buf, r.rawRequestLine[idx+2:]...)
		r.rawRequestLine = r.rawRequestLine[:idx]
		r.state = requestStateParsingHeaders

		r.Line, err = parseRawRequestLine(r.rawRequestLine[:idx])
		if err != nil {
			return n, err
		}
	}

	return n, nil
}

func parseRawRequestLine(buf []byte) (Line, error) {
	split := bytes.Split(buf, []byte(" "))

	if len(split) != 3 {
		return Line{}, fmt.Errorf("invalid request line")
	}

	method := split[0]
	if len(method) == 0 {
		return Line{}, fmt.Errorf("no method")
	}

	target := split[1]
	if len(target) == 0 {
		return Line{}, fmt.Errorf("no target")
	}

	protocol := split[2]
	httpVersion := bytes.Split(protocol, []byte("/"))
	if len(httpVersion) != 2 {
		return Line{}, fmt.Errorf("invalid http version")
	}

	if !bytes.Equal(httpVersion[0], []byte("HTTP")) {
		return Line{}, fmt.Errorf("not http")
	}

	version := httpVersion[1]

	return Line{
		Method:  string(method),
		Target:  string(target),
		Version: string(version),
	}, nil
}

func (r *Request) parseHeaders(reader io.Reader) (int, error) {
	buf := r.currLine

	n, err := reader.Read(buf)
	if err != nil {
		return 0, err
	}

	buf = buf[:n]

	r.rawHeaders = append(r.rawHeaders, buf...)

	idx := bytes.Index(r.rawHeaders, []byte("\r\n\r\n"))
	if idx == -1 {
		return n, nil
	}

	var headersN int

	for {
		m, done, err := r.Headers.Parse(r.rawHeaders[headersN : idx+4])
		headersN += m
		if err != nil {
			return 0, err
		}

		if done {
			// if there is something left in the buffer, append it to the request body
			r.Body = append(r.Body, r.rawHeaders[idx+4:]...)
			r.state = parsinBody
			return n, nil
		}
	}
}
