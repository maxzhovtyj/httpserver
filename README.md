# httpfromtcp

`httpfromtcp` is a learning-focused Go project that builds HTTP functionality directly on top of raw TCP connections.  
The project includes:
```markdown
- A minimal HTTP request parser
- A basic HTTP response writer
- A small TCP-based HTTP server
- Header parsing utilities
- Chunked transfer response examples
- Trailer support examples
- Simple TCP and UDP networking experiments
```

The goal of the project is to understand how HTTP works underneath the usual `net/http` abstractions.

## Project Structure
```text
httpfromtcp/
├── cmd/
│   ├── httpserver/
│   │   ├── main.go
│   │   └── vim-vs-neovim-prime.mp4
│   ├── tcplistener/
│   │   └── main.go
│   └── udpsender/
│       └── main.go
├── internal/
│   ├── headers/
│   │   ├── headers.go
│   │   └── headers_test.go
│   ├── request/
│   │   ├── request.go
│   │   └── request_test.go
│   ├── response/
│   │   └── response.go
│   └── server/
│       ├── handler.go
│       └── server.go
├── go.mod
├── messages.txt
└── .gitignore
```