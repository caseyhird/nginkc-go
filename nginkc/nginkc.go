package nginkc

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const (
	PORT = 8080
	HOST = "127.0.0.1"
)

type Request struct {
	Method   string
	Path     string
	Protocol string
	Headers  map[string]string
	Body     string
}

func (req Request) String() string {
	return fmt.Sprintf("Request{Method: %s, Path: %s, Protocol: %s, Headers: %v, Body: %s}", req.Method, req.Path, req.Protocol, req.Headers, req.Body)
}

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

func (res Response) toBytes() []byte {
	response_str := fmt.Sprintf("HTTP/1.1 %d %s\n", res.StatusCode, http.StatusText(res.StatusCode))
	for key, value := range res.Headers {
		response_str += fmt.Sprintf("%s: %s\n", key, value)
	}
	response_str += "\n" + res.Body
	return []byte(response_str)
}

type App interface {
	Call(req Request) Response
}

func parseRequest(conn net.Conn) (Request, error) {
	request := Request{}
	reader := bufio.NewReader(conn)

	// Get first line
	first_line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading first line:", err)
		return request, err
	}
	first_line_split := strings.Split(first_line, " ")
	request.Method = strings.TrimSpace(first_line_split[0])
	request.Path = strings.TrimSpace(first_line_split[1])
	request.Protocol = strings.TrimSpace(first_line_split[2])

	if request.Protocol != "HTTP/1.1" {
		fmt.Println("Unsupported protocol:", request.Protocol)
		return request, fmt.Errorf("unsupported protocol: %s", request.Protocol)
	}

	// Get headers until empty line
	request.Headers = make(map[string]string)
	for {
		header_line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading headers:", err)
			return request, err
		}
		if header_line == "\r\n" {
			break
		}
		header_split := strings.Split(header_line, ": ")
		request.Headers[header_split[0]] = strings.TrimSpace(header_split[1])
	}

	// All remaining data is body
	// NOTE: Need to read using "Content-Length" header. Can't read until reaching EOF
	// because http stays "alive" and doesn't close the connection
	content_length_str, ok := request.Headers["Content-Length"]
	if !ok {
		fmt.Println("Content-Length header not found. Assuming no body.")
		return request, nil
	}
	content_length, err := strconv.Atoi(content_length_str)
	if err != nil {
		fmt.Println("Error converting Content-Length to int:", err)
		return request, err
	}
	buff := make([]byte, content_length)
	body_data, err := io.ReadFull(reader, buff)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return request, err
	}
	request.Body = string(buff[:body_data])
	return request, nil
}

func handleConnection(conn net.Conn, app App) {
	defer conn.Close()

	request, err := parseRequest(conn)
	if err != nil {
		fmt.Println("Error parsing request:", err)
		return
	}

	response := app.Call(request)
	conn.Write(response.toBytes())
}

func Serve(app App) {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", HOST, PORT))
	if err != nil {
		fmt.Println("Error binding socket:", err)
		return
	}
	defer listener.Close()
	fmt.Printf("Serving app on %s:%d\n", HOST, PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn, app)
	}
}
