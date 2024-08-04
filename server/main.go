package main

import (
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
		conn, err := s.l.Accept()
		if err != nil {
			s.log.Error("error accepting connection: %s", err.Error())
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		s.log.Infof("closing connection on %s", conn.RemoteAddr())
		conn.Close()
	}()
	s.log.Infof("new connection on %s", conn.RemoteAddr())
	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Println("Error reading data: ", err.Error())
		}

		msg := make([]byte, n)
		copy(msg, data[:n])

		log.Println("read:", string(msg))

		parseCommand(string(msg))

		conn.Write([]byte("+OK\r\n"))
	}
}

func (s *Server) Start() error {
	l, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}

	s.log.Infof("accepting connections on %s", s.ListenAddr)

	s.l = l
	return nil
}

func NewServer(config Config, log *zap.SugaredLogger) *Server {
	return &Server{
		Config: config,
		log:    log,
	}
}

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	s := NewServer(Config{ListenAddr: "0.0.0.0:36379"}, logger.Sugar())
	err = s.Start()
	if err != nil {
		log.Fatal(err)
	}

	s.acceptLoop()
}
