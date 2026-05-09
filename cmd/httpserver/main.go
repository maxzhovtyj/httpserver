package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"httpfromtcp/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const port = 42069

func main() {
	srv := server.New(port)

	srv.HandleFunc("/", func(req *request.Request, resp *response.Response) {
		resp.Headers.Set("Content-Type", "text/html")
		resp.Write([]byte(`
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
`))
	})
	srv.HandleFunc("/ping", func(req *request.Request, resp *response.Response) {
		resp.Write([]byte("pong"))
	})
	srv.HandleFunc("/your-problem", func(req *request.Request, resp *response.Response) {
		resp.StatusCode = response.StatusBadRequest
		resp.Headers.Set("Content-Type", "text/html")
		resp.Write([]byte(`
<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
`))
	})
	srv.HandleFunc("/my-problem", func(req *request.Request, resp *response.Response) {
		resp.StatusCode = response.StatusInternalServerError
		resp.Headers.Set("Content-Type", "text/html")
		resp.Write([]byte(`
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
`))
	})
	srv.HandleFunc("/httpbin/100", stream)
	srv.HandleFunc("/trailers", trailers)
	srv.HandleFunc("/video", video)

	err := srv.Serve()
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	defer srv.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func stream(req *request.Request, resp *response.Response) {
	st, err := http.Get("https://httpbin.org/stream/100")
	if err != nil {
		return
	}

	defer st.Body.Close()

	resp.StatusCode = response.StatusOK
	resp.Headers.Set("Transfer-Encoding", "chunked")
	resp.Flush()

	buf := make([]byte, 1024)

	for {
		n, err := st.Body.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return
		}

		buf = buf[:n]

		resp.WriteChunked([]byte(fmt.Sprintf("%x\r\n", len(buf))))
		resp.WriteChunked(buf)
		resp.WriteChunked([]byte("\r\n"))
	}

	resp.WriteChunked([]byte(fmt.Sprintf("0\r\n")))
	resp.WriteChunked([]byte("\r\n"))
}

func trailers(req *request.Request, resp *response.Response) {
	st, err := http.Get("https://httpbin.org/html")
	if err != nil {
		return
	}

	defer st.Body.Close()

	resp.Headers.Set("Transfer-Encoding", "chunked")
	resp.Headers.Set("Trailers", "x-content-sha256, x-content-length")
	resp.StatusCode = response.StatusOK

	resp.Flush()

	buf := make([]byte, 1024)

	h := sha256.New()
	totalLen := 0

	for {
		n, err := st.Body.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return
		}

		buf = buf[:n]

		totalLen += len(buf)
		_, _ = h.Write(buf)

		resp.WriteChunked([]byte(fmt.Sprintf("%x\r\n", len(buf))))
		resp.WriteChunked(buf)
		resp.WriteChunked([]byte("\r\n"))
	}

	resp.WriteChunked([]byte(fmt.Sprintf("0\r\n")))
	resp.Trailers.Set("X-Content-SHA256", string(h.Sum(nil)))
	resp.Trailers.Set("X-Content-Length", strconv.Itoa(totalLen))
	resp.FlushTrailers()
	resp.WriteChunked([]byte("\r\n"))
}

func video(req *request.Request, resp *response.Response) {
	f, err := os.Open("/Users/maksymzhovtaniuk/Desktop/Programming/Go/httpfromtcp/cmd/httpserver/vim-vs-neovim-prime.mp4")
	if err != nil {
		return
	}
	defer f.Close()

	resp.SetStatusCode(200)
	resp.Headers.Set("Content-Type", "video/mp4")
	resp.Headers.Set("Transfer-Encoding", "chunked")
	resp.Flush()

	buf := make([]byte, 1024*1000)

	for {
		n, err := f.Read(buf)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		}

		buf = buf[:n]

		resp.WriteChunked([]byte(fmt.Sprintf("%x\r\n", len(buf))))
		resp.WriteChunked(buf)
		resp.WriteChunked([]byte("\r\n"))
	}

	resp.WriteChunked([]byte(fmt.Sprintf("0\r\n")))
	resp.WriteChunked([]byte("\r\n"))
}
