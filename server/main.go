package main

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"net"
)

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	l   net.Listener
	log *zap.SugaredLogger
}

func (s *Server) acceptLoop() {
	for {
		s.log.Infof("accepting connections on %s", s.ListenAddr)
		conn, err := s.l.Accept()
		if err != nil {
			log.Fatalf("error accepting connection: %s", err.Error())
		}
		go s.process(conn)
	}
}

func (s *Server) process(conn net.Conn) {
	defer func() {
		s.log.Infof("closing connection on %s", conn.RemoteAddr())
		conn.Close()
	}()

	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error reading data: ", err.Error())
		}

		log.Printf("Received %d bytes: %s\n", n, data[:n])

		resp := fmt.Sprint("+PONG\r\n")

		_, err = conn.Write([]byte(resp))
		if err != nil {
			log.Println("Error writing data: ", err.Error())
		}
	}
}

func NewServer(config Config, log *zap.SugaredLogger) (*Server, error) {
	l, err := net.Listen("tcp", config.ListenAddr)
	if err != nil {
		return nil, err
	}

	return &Server{
		Config: config,
		l:      l,
		log:    log,
	}, nil
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	s, err := NewServer(Config{ListenAddr: "0.0.0.0:6379"}, logger.Sugar())
	if err != nil {
		log.Fatal(err)
	}

	s.acceptLoop()
}
